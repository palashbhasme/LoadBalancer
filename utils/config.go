package utils

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/palashbhasme/loadbalancer/internals"
	"gopkg.in/yaml.v3"
)

// loadFile reads the file from the given path
func loadFile(path *string) ([]byte, error) {

	info, err := os.Stat(*path)

	if os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	} else if info.Size() == 0 {
		return nil, errors.New("file is empty")
	}

	data, err := os.ReadFile(*path)
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}

	log.Printf("data: %s", string(data))

	return data, nil

}

// LoadConfig reads the config file and unmarshals it into a Servers struct
func LoadConfig() (*internals.Servers, error) {

	var config internals.Servers

	// Set a new mutex to ensure thread-safe operations on the Servers struct
	config.SetMu(&sync.Mutex{})
	config.SetIndex(0)
	path := flag.String("path", "", "path to config file")
	flag.Parse()

	if *path == "" {
		return nil, errors.New("path to config file is required")
	}

	data, err := loadFile(path)
	if err != nil {
		return nil, errors.New("error loading config file: " + err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New("error unmarshalling YAML data into Servers struct: " + err.Error())
	}
	log.Printf("config: %+v", config)
	return &config, nil
}
