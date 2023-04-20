/* Package for creating and managing the config file.
 */

package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Name string
	OS   string
}

const configFileName = "config.json"
const configFolderName = "anio"

// [Load]: Checks user's OS, then reads config data from file or creates new.
func Load() (Config, error) {
	var configFolderPath string
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		configFolderPath, _ = os.UserConfigDir()
		configFolderPath = filepath.Join(configFolderPath, configFolderName)
	default:
		log.Fatal().Msgf("unsupported OS: %s", runtime.GOOS)
	}

	configFilePath := filepath.Join(configFolderPath, configFileName)

	if !exists(configFilePath) {
		log.Info().Msg("Config file not found. Creating a default one.")
		conf, err := createDefaultConfig(configFolderPath)
		if err != nil {
			log.Fatal().Err(err).Msg("Couldn't create config.")
		}
		return conf, nil
	}

	log.Info().Msg("Found existing config file.")
	conf, err := loadExistingConfig(configFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't load config from file.")
	}
	return conf, nil
}

// check file existence
func exists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Load config from file.
func loadExistingConfig(filePath string) (Config, error) {
	conf := Config{}
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read file.")
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't load data from file to memory.")
	}
	return conf, err
}

// Create a new config, write to file and load it.
func createDefaultConfig(folderPath string) (Config, error) {
	// check if directory exists just in case:
	if !exists(folderPath) {
		_ = os.Mkdir(folderPath, os.FileMode(0777)) // permissions for linux
	}

	// create the conf
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't read user personal data.")
	}
	conf := Config{
		OS:   runtime.GOOS,
		Name: currentUser.Username,
	}

	// convert to json
	configData, _ := json.Marshal(conf)

	// write it to file
	fileName := filepath.Join(folderPath, configFileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't create file.")
	}

	defer file.Close()
	_, err = file.Write(configData)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't write data to file.")
	}

	return conf, err
}
