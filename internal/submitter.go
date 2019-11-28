package internal

import (
	"Ansem/internal/submitters"
	"context"
	"log"
	"sync"
)

func StartSubmitter(submitterCtx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	toSubmit := submitterCtx.Value("submit").(<-chan string)
	//Init submitters
	submitFunction := submitters.RuCTFSubmitHTTP

	//Create a thread safe map to verify flags
	var submitted sync.Map

	//Create channel to pass filtered flags
	flagChannel := make(chan string, 10)
	submitterCtx = context.WithValue(submitterCtx, "flagChannel", flagChannel)
	submitterCtx = context.WithValue(submitterCtx, "alreadySubmitted", &submitted)

	//Start the submitter
	go submitFunction(submitterCtx)
	//Check if the flags are already submitted
	for flag := range toSubmit {
		//The regex is checked via exploiter
		//If is present or doesn't match the flag regexp continue
		if _, result := submitted.Load(flag); result {
			log.Printf("SUBMITTER:\t flag %s already retrieved!\n", flag)
			continue
		} else {
			flagChannel <- flag
			log.Printf("SUBMITTER:\t flag %s is new, now i will send it!\n", flag)
		}
	}

}
