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

// config data container to be used in the main package.
type Config struct {
	Name string
	OS   string
}

// Checks User's OS, then looks for the config file and tries to read the
// data.

// If there is no config file, creates a default one.
func Load() Config {

	//question about configpath: Is it better to split it into path-filename,
	//or it doesn't really matter?

	//I initially wanted to split them, but got confused by string
	//concatenation in Go and decided to opt into a single filepath+filename
	//string for now.
	var configPath string
	switch runtime.GOOS {
	case "windows":
		configPath = "./config.json"
	default:
		log.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	conf := Config{} //TODO:I guess I can just delete this now?

	// TODO: Maybe this should be the other way around now?
	// Load existing conf by default, create one if not found/errored
	// on loading.
	if exists(configPath) {
		fmt.Println("Found existing config file.")
		// TODO: should I catch an error// make these funcs return possible
		// errors too?
		conf = loadExistingConfig(configPath)
		// TODO: should check if everything is loaded successfully.
	}
	fmt.Println("Config file not found. Creating a default one.")

	conf = createDefaultConfig(configPath)

	// minor thing: I don't really like "conf" variable name but I'm stumped
	// how else to name it comprehensibly. Any ideas?
	return conf
}

// Another question: is storing these functions outside Load() scope fine or
// should have I put them inside Load()? I kinda get confused about
// these things.

// check file existence
func exists(filepath string) bool {

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

// Just to make sure: is it okay to abstract parts of code away into
// separate fns like these? I like doing this because it makes it easier to
// understand the overall flow, but I'm concerned if it unnecessarily
// spagettifies the code?

// Load config from file.
func loadExistingConfig(filepath string) Config {
	conf := Config{}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf //err?
}

// Create a new config, write to file and load it.
func createDefaultConfig(filepath string) Config {

	//Just making sure: Are these three comment lines fine?
	//They kinda help to divide the code visually into functional bits.

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

	return conf //err?
}
