package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type conf struct {
	Directory  string `yaml:"exploits_dir"`
	Tick       int    `yaml:"tick"`
	GameServer string `yaml:"gameserver"`
	Workers    int    `yaml:"workers"`
}

func main() {
	var c conf
	c.getConf()
	toSubmit := make(chan string, 20)
	teams := make(chan string, 20)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		StartExploiter(c.Directory, toSubmit, teams, c.Tick, c.Workers)
	}()
	go func() {
		defer wg.Done()
		StartSubmitter(c.GameServer, toSubmit)
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
