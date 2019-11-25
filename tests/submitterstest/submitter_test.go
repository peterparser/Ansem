package submitters

import (
	"Ansem/internal/submitters"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

type RuctfeServer struct{}

const gameserver = "127.0.0.1:8080"

func (s *RuctfeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Implement testing here
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message"": "NOT FOUND"}`))
		return
	}
	if r.Header.Get("X-Team-Token") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message"": "NOT FOUND"}`))
		return

	}
	decoder := json.NewDecoder(r.Body)
	var t []string
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message"": "NOT FOUND"}`))
		return

	}
	if len(t) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message"": "NOT FOUND"}`))
		return

	}
	var response []submitters.RuCtfFlag
	for _, flag := range t {
		response = append(response, submitters.RuCtfFlag{"ok", flag, true})
	}
	marshalled, _ := json.Marshal(&response)
	w.WriteHeader(http.StatusOK)
	w.Write(marshalled)

}

func TestRuctfeSubmit(t *testing.T) {

	// Start http server
	go func() {
		s := &RuctfeServer{}
		http.Handle("/", s)
		log.Fatal(http.ListenAndServe(gameserver, nil))
	}()

	time.Sleep(10 * time.Millisecond)

	var flags []string
	flags = append(flags, "1")
	flags = append(flags, "2")
	flags = append(flags, "3")

	flagJson, err := json.Marshal(&flags)
	if err != nil {
		t.Errorf("Error in json marshalling %v\n", err)
	}
	req, err := http.NewRequest("PUT", "http://"+gameserver, bytes.NewBuffer(flagJson))
	if err != nil {
		t.Errorf("Error in creating request struct %v\n", err)
	}
	//Add headers
	req.Header.Set("X-Team-Token", "token")

	//Send flag
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in executing request %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Error in request status should be %d instead of %d\n", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()
	var flagResult []submitters.RuCtfFlag
	//Parse response

	err = json.NewDecoder(resp.Body).Decode(&flagResult)
	if err != nil {
		t.Errorf("Error in json decoding %v\n", err)
	}

	for _, res := range flagResult {
		if !res.Status {
			t.Errorf("Status should be %t instead of %t", true, res.Status)
		}
	}

}
