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
	var submitHandler = make(map[string]func(string, string, <-chan string, *sync.Map, string))
	//Define submission method
	submitHandler["TCP"] = submitters.RuCTFSubmitNC
	submitHandler["HTTP"] = submitters.RuCTFSubmitHTTP

	//Create a thread safe map to verify flags
	var submitted sync.Map

	//Create channel to pass filtered flags
	flagChannel := make(chan string, 10)

	//Start the submitter
	go submitHandler[subType](gameServer, flagAccepted, flagChannel, &submitted, token)
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
