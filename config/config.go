package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
)

func Load() Config {

	// Checks User's OS, then looks for the config file and tries to read the
	// data.
	// If there is no config file, creates a deafult one.
	// Returns Config struct that can be used in main package.

	currentOs := runtime.GOOS

	var configPath string
	switch currentOs {
	case "windows":
		configPath = "./config.ini"
	default:
		configPath = "./config.ini" // placeholder?
	}

	conf := Config{}
	if exists(configPath) {
		fmt.Println("Found existing config file.")
		conf = loadExistingConfig(configPath)
		// should check if everything is loaded successfully, but I dunno
		// how.
	} else {
		fmt.Println("Config file not found. Creating a default one.")
		conf = createDefaultConfig(configPath)
	}
	return conf
}

func exists(filepath string) bool {
	//this just returns true or false if file already exists at the provided
	//path.

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

type Config struct {
	// an object that gets returned to the main package, containing all the
	// needed parameters as fields.
	Name string
	Os   string
}

func loadExistingConfig(filepath string) Config {
	// loads data from json at the provided filepath, creates a config struct
	// with it
	conf := Config{}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &conf)
	return conf
}

func createDefaultConfig(filepath string) Config {
	//Creates a new config with the default data: username and os used.
	//Then saves it at the default folder.

	//create the conf
	conf := Config{}
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	conf.Os = runtime.GOOS
	conf.Name = currentUser.Username

	//convert to json
	configData, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	//write it to file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(configData))

	return conf
}
