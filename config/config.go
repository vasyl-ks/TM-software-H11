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

type client struct {
	FileDir 	string `json:"fileDir"`
}

var Server server
var Client client

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}
	defer file.Close()

	temp := struct {
		S server `json:"server"`
		C client `json:"client"`
	}{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&temp)
	if err != nil {
		fmt.Println("Error loading configuration: ", err)
		return
	}

	Server = temp.S
	Client = temp.C

	Server.Interval = time.Duration(Server.I) * time.Millisecond
}