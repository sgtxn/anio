/* Package for creating and managing the config file.
 */

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"anio/pkg/duration"
	"anio/providers/anilist/consts"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Name          string         `json:"name"`
	OS            string         `json:"os"`
	Inputs        *InputsConfig  `json:"inputs,omitempty"`
	AnilistConfig *AnilistConfig `json:"anilist,omitempty"`

	lock     sync.Mutex
	selfPath string
}

const (
	configFileName = "config.json"
)

func GetConfigFilePath(projectFolderName string) (string, error) {
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		configFolderPath, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("cannot access user config folder: %w", err)
		}
		return filepath.Join(configFolderPath, projectFolderName, configFileName), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// SaveConfig updates the locally saved config.
func (cfg *Config) SaveConfig() error {
	cfg.lock.Lock()
	defer cfg.lock.Unlock()

	if err := writeConfigToFile(cfg); err != nil {
		return fmt.Errorf("couldn't update config file under %s: %w", cfg.selfPath, err)
	}

	return nil
}

// Load checks user's OS, then reads config data from file or creates a new default config.
func Load(configFilePath string) (*Config, error) {
	if !exists(configFilePath) {
		log.Info().Msg("config file not found. Creating a default one.")
		conf, err := createDefaultConfig(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("couldn't create config: %w", err)
		}

		return conf, nil
	}

	log.Info().Msg("found existing config file")
	conf, err := loadExistingConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config from file: %w", err)
	}

	return conf, nil
}

func loadExistingConfig(cfgFilePath string) (*Config, error) {
	log.Info().Msgf("loading config from %s...", cfgFilePath)

	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %w", err)
	}

	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("couldn't load data from file to memory: %w", err)
	}

	conf.selfPath = cfgFilePath

	return &conf, nil
}

func createDefaultConfig(cfgFilePath string) (*Config, error) {
	cfgFolderPath := filepath.Dir(cfgFilePath)
	log.Info().Msgf("creating default config under %s...", cfgFolderPath)

	if !exists(cfgFolderPath) {
		err := os.Mkdir(cfgFolderPath, 0o744)
		if err != nil {
			return nil, fmt.Errorf("couldn't create folder: %w", err)
		}
	}

	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("couldn't read user personal data: %w", err)
	}

	cfg := getDefaultConfig(runtime.GOOS, currentUser.Username)
	cfg.selfPath = filepath.Join(cfgFolderPath, configFileName)

	err = writeConfigToFile(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't write config file: %w", err)
	}

	return cfg, nil
}

func writeConfigToFile(cfg *Config) error {
	cfgBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	file, err := os.Create(cfg.selfPath)
	if err != nil {
		return fmt.Errorf("couldn't create file %s: %w", cfg.selfPath, err)
	}

	defer file.Close()

	if _, err = file.Write(cfgBytes); err != nil {
		return fmt.Errorf("couldn't write data to file %s: %w", cfg.selfPath, err)
	}

	if err = file.Chmod(0o644); err != nil {
		return fmt.Errorf("couldn't chmod config file %s: %w", cfg.selfPath, err)
	}

	return nil
}

func exists(cfgFilePath string) bool {
	_, err := os.Stat(cfgFilePath)
	return !os.IsNotExist(err)
}

func getDefaultConfig(userOS, username string) *Config {
	return &Config{
		OS:   userOS,
		Name: username,
		Inputs: &InputsConfig{
			LocalPollers: &LocalAppConfig{
				PollingInterval: duration.Duration{Duration: time.Second * 5}, //nolint:gomnd // it's fine to hardcode the defaults
				MpvConfig: &MpvConfig{
					Enabled:       false,
					UseJSONRPCAPI: false,
				},
			},
		},
		AnilistConfig: &AnilistConfig{
			Auth: AnilistAuthConfig{
				ClientID:     consts.ClientID,
				ClientSecret: consts.ClientSecret,
			},
		},
	}
}
