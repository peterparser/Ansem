package internal

import (
	"bufio"
	"log"
	"os"
)

func GetTeamAsChan(fileTeam string) chan string {
	teamChannel := make(chan string, 20)
	file, err := os.Open(fileTeam)
	if err != nil {
		log.Fatalf("TEAM\t Error Team File:\t  %s \n", fileTeam)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	go func() {
		for _, line := range lines {
			teamChannel <- line
		}
	}()
	return teamChannel
}
