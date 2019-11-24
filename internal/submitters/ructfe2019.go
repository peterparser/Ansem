package submitters

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RuCtfFlag struct {
	Msg    string `json:"msg"`
	Flag   string `json:"flag"`
	Status bool   `json:"status"`
}

func RuCTFSubmitHTTP(gameServer string, acceptedFlag string, flagChannel <-chan string, alreadySubmitted *sync.Map, token string) {
	//	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//Create the tcp connection
	var flags []string
	for {
		select {
		//Read the flag
		case flag := <-flagChannel:
			flags = append(flags, flag)
		//Create json from flag
		case <-time.After(5 * time.Second):
			if flags != nil {
				flagJson, err := json.Marshal(flags)
				if err != nil {

					log.Fatalf("SUBMITTER\nError in json marshal with %s\nTrace: %s\n", gameServer, err)
				}
				req, err := http.NewRequest("PUT", gameServer, bytes.NewBuffer(flagJson))
				//Add headers
				req.Header.Set("X-Team-Token", token)
				if err != nil {
					log.Fatalf("SUBMITTER\tConnection Error HTTP:\t Server %s\n Trace:%s\n", gameServer, err)
				}
				//Send flag
				client := &http.Client{
					Timeout: time.Second * 5,
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("SUBMITTER\tError Send Flag:\t Server %s\nTrace: %s\n", gameServer, err)
				}
				defer resp.Body.Close()
				var flagResult []RuCtfFlag
				//Parse response

				err = json.NewDecoder(resp.Body).Decode(&flagResult)
				if err != nil {
					log.Fatalf("SUBMITTER\tError Unmarshalling Flag:\nTrace: %s\n", err)
				}
				for _, flagStatus := range flagResult {
					if flagStatus.Status {
						alreadySubmitted.Store(flagStatus.Flag, true)
					} else {
						log.Printf("SUBMITTER\tInvalid Flag:\t %s \n", flagStatus.Flag)
					}
				}
				flags = nil
			}
		}
	}

}

func RuCTFSubmitNC(gameServer string, acceptedFlag string, flagChannel <-chan string, alreadySubmitted *sync.Map, token string) {

	//Create the tcp connection
	connection, err := net.DialTimeout("tcp", gameServer, 10*time.Second)
	if err != nil {
		log.Fatalf("SUBMITTER\tConnection Error TCP:\t Server %s\n Trace:%s\n", gameServer, err)
	}
	for {
		//Buffered reader
		reader := bufio.NewReader(connection)
		select {
		//Read the flag
		case flag := <-flagChannel:
			//Send the flag
			fmt.Fprintf(connection, "%s\n", flag)
			//Read the response
			response, _ := reader.ReadString('\n')
			//If it's accepted, store it
			if strings.Contains(response, acceptedFlag) {
				alreadySubmitted.Store(flag, true)

			}
			//After x seconds without flag, stop
		case <-time.After(10 * time.Second):
			connection.Close()
			fmt.Print("Chiudo\n")
			time.Sleep(10 * time.Second)
			connection, err = net.DialTimeout("tcp", gameServer, 10*time.Second)
			if err != nil {
				log.Fatalf("SUBMITTER\tConnection Error TCP:\t Server %s\n Trace:%s\n", gameServer, err)
			}
		}
	}
}
