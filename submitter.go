package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"net/http"
	"crypto/tls"
	"bytes"
	"io/ioutil"
)

func jsonRuCTF(flag string)[]byte{
	encode:=fmt.Sprintf("[\"%s\"]}",flag)
	return []byte(encode)
}





func StartSubmitter(submitterCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	gameServer := submitterCtx.Value("gameServer").(string)
	toSubmit := submitterCtx.Value("submit").(chan string)
	subType := submitterCtx.Value("subType").(string)
	flagAccepted := submitterCtx.Value("flagAccepted").(string)
	token := submitterCtx.Value("token").(string)
	//flagRegex, err := regexp.Compile(submitterCtx.Value("flagRegex").(string))

	//if err != nil {
	//	log.Fatalf("Invalid regexp\n")
	//}

	//Init submitters
	var submitHandler = make(map[string]func(string, string, <-chan string, chan<- string, string))
	submitHandler["TCP"] = submitNC
	submitHandler["HTTP"] = submitHTTP
	//Define submission method

	//Create a map to verify flags
	submitted := make(map[string]bool)
	//Create channel to pass filtered flags
	flagChannel := make(chan string, 10)
	//Create the channel to communicate with the map handler
	mapWrite := make(chan string)
	mapRead := make(chan string)
	mapGet := make(chan bool)

	//Start the handler of the map
	go func() {
		for {
			select {
			case write := <-mapWrite:
				submitted[write] = true
			case read := <-mapRead:
				_, found := submitted[read]
				mapGet <- found
			}
		}
	}()

	//Start the submitter
	go submitHandler[subType](gameServer, flagAccepted, flagChannel, mapWrite, token )
	//Check if the flags are already submitted
	for flag := range toSubmit {
		mapRead <- flag
		present := <-mapGet
		//The regex is checked via exploiter
		//If is present or doesn't match the flag regexp continue
		if present { //||  !matched {
			continue
		} else {
			flagChannel <- flag
		}
	}

}

func submitHTTP(gameServer string, acceptedFlag string, flagChannel <-chan string, handler chan<- string, token string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//Create the tcp connection
	for {
		//Buffered reader
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			flagJson := jsonRuCTF(flag)
			req, err := http.NewRequest("PUT", gameServer, bytes.NewBuffer(flagJson))
			req.Header.Set("X-Team-Token",token)
			if err != nil {
				log.Fatalf("Error in connection http  with %s %s\n", gameServer,err)
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{
				Timeout: time.Second * 5,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error in sending flag with %s %s\n", gameServer,err)
			}
			fmt.Printf("Sent: %s\n", flag)
			fmt.Printf("response status %s",resp.Status)
			response, err := ioutil.ReadAll(resp.Body)
			if err!= nil {
				log.Fatalf("Error in receving answer with %s %s\n", gameServer,err)
			}
			if strings.Contains(string(response), acceptedFlag) {
				handler <- flag

			}
		}
	}

}


func submitNC(gameServer string, acceptedFlag string, flagChannel <-chan string, handler chan<- string, token string) {

	//Create the tcp connection
	connection, err := net.DialTimeout("tcp", gameServer, 10*time.Second)
	if err != nil {
		log.Fatalf("Error in connection with %s %s\n", gameServer,err)
	}
	for {
		//Buffered reader
		reader := bufio.NewReader(connection)
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			fmt.Printf("Sent: %s\n", flag)
			fmt.Fprintf(connection, "%s\n", flag)
			//Read the response
			response, _ := reader.ReadString('\n')
			//If it's accepted, store it
			if strings.Contains(response, acceptedFlag) {
				handler <- flag

			}
			//After x seconds without flag, stop
		case <-time.After(10 * time.Second):
			connection.Close()
			fmt.Print("Chiudo\n")
			time.Sleep(10 * time.Second)
			connection, err = net.DialTimeout("tcp", gameServer, 10*time.Second)
			if err != nil {
				log.Fatalf("Error in connection with %s\n", gameServer)
			}
		}
	}
}
