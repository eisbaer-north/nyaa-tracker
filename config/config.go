package config

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Path string
	Autostart bool
}

func LoadConfig(path string) Config {
	//Open the config file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	//Create a variable for the configuration and load the json data into
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
