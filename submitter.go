package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

func StartSubmitter(submitterCtx context.Context) {

	gameServer := submitterCtx.Value("gameServer").(string)
	toSubmit := submitterCtx.Value("submit").(chan string)
	subType := submitterCtx.Value("subType").(string)
	flagAccepted := submitterCtx.Value("flagAccepted").(string)
	flagRegex, err := regexp.Compile(submitterCtx.Value("flagRegex").(string))

	if err != nil {
		log.Fatalf("Invalid regexp\n")
	}

	//Init submitters
	var submitHandler = make(map[string]func(string, string, <-chan string, chan<- string))
	submitHandler["TCP"] = submitNC
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
	go submitHandler[subType](gameServer, flagAccepted, flagChannel, mapWrite)
	//Check if the flags are already submitted
	for flag := range toSubmit {
		mapRead <- flag
		present := <-mapGet
		matched := flagRegex.MatchString(flag)
		//If is present or doesn't match the flag regexp continue
		if present || !matched {
			continue
		} else {
			flagChannel <- flag
		}
	}

}

func submitNC(gameServer string, acceptedFlag string, flagChannel <-chan string, handler chan<- string) {

	//Create the tcp connection
	connection, err := net.DialTimeout("tcp", gameServer, 10*time.Second)
	if err != nil {
		log.Fatalf("Error in connection with %s\n", gameServer)
	}
	for {
		//Buffered reader
		reader := bufio.NewReader(connection)
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			fmt.Printf("Sto per inviare %s\n", flag)
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
