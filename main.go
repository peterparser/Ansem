package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"text/tabwriter"
)

type conf struct {
	Directory      string `yaml:"exploits_dir"`
	Tick           int    `yaml:"tick"`
	GameServer     string `yaml:"gameserver"`
	Workers        int    `yaml:"workers"`
	TeamFile       string `yaml:"team_file"`
	SubmissionType string `yaml:"submission_type"`
}

func main() {
	var c conf
	c.getConf()

	//Fix path if the last char is not "/"
	if c.Directory[len(c.Directory)-1] != '/' {
		c.Directory = fmt.Sprintf("%s/", c.Directory)
	}

	//Aligned print
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, _ = fmt.Fprintf(writer, "Hi, I'm starting with these settings:\n\nExploits Dir:\t%s\nGameserver:\t%s\nTeamfile:\t%s\nTick:\t%d\nWorkers:\t%d\n",
		c.Directory, c.GameServer, c.TeamFile, c.Tick, c.Workers)
	writer.Flush()

	toSubmit := make(chan string, 20)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		StartExploiter(c.Directory, toSubmit, c.TeamFile, c.Tick, c.Workers)
	}()
	go func() {
		defer wg.Done()
		StartSubmitter(c.GameServer, toSubmit, c.SubmissionType)
	}()
	wg.Wait()
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
