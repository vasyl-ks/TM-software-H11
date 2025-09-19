package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type server struct {
	ClientPort 	int	`json:"clientPort"`
	I 			int	`json:"intervalMiliSeconds"`
	Interval 	time.Duration
}

var Server server

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Server)
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}

	Server.Interval = time.Duration(Server.I) * time.Millisecond
}