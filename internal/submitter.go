package internal

import (
	"Ansem/internal/submitters"
	"context"
	"log"
	"sync"
)

func StartSubmitter(submitterCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	gameServer := submitterCtx.Value("gameServer").(string)
	toSubmit := submitterCtx.Value("submit").(chan string)
	subType := submitterCtx.Value("subType").(string)
	flagAccepted := submitterCtx.Value("flagAccepted").(string)
	token := submitterCtx.Value("token").(string)

	//Init submitters
	var submitHandler = make(map[string]func(string, string, <-chan string, chan<- string, string))
	//Define submission method
	submitHandler["TCP"] = submitters.RuCTFSubmitNC
	submitHandler["HTTP"] = submitters.RuCTFSubmitHTTP

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
	go submitHandler[subType](gameServer, flagAccepted, flagChannel, mapWrite, token)
	//Check if the flags are already submitted
	for flag := range toSubmit {
		mapRead <- flag
		present := <-mapGet
		//The regex is checked via exploiter
		//If is present or doesn't match the flag regexp continue
		if present {
			log.Printf("SUBMITTER:\t flag %s already retrieved!\n", flag)
			continue
		} else {
			flagChannel <- flag
			log.Printf("SUBMITTER:\t flag %s is new, now i will send it!\n", flag)
		}
	}

}
