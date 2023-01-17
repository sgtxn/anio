/* Package for creating and managing the config file.
 */

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
)

type Config struct {
	Name string
	OS   string
}

// [Load]: Checks user's OS, then reads config data from file or creates new.
func Load() (Config, error) {

	var configPath string
	switch runtime.GOOS {
	case "windows":
		configPath = "./config.json"
	default:
		log.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	if !exists(configPath) {
		fmt.Println("Config file not found. Creating a default one.")
		conf, err := createDefaultConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}
		return conf, nil
	}

	fmt.Println("Found existing config file.")
	conf, err := loadExistingConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return conf, nil
}

// check file existence
func exists(filepath string) bool {

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

// Load config from file.
func loadExistingConfig(filepath string) (Config, error) {
	conf := Config{}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf, nil
}

// Create a new config, write to file and load it.
func createDefaultConfig(filepath string) (Config, error) {

	//create the conf
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	conf := Config{
		OS:   runtime.GOOS,
		Name: currentUser.Username,
	}

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
	_, err = file.Write([]byte(configData))
	if err != nil {
		log.Fatal(err)
	}

	return conf, nil
}
