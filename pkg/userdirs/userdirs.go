package userdirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	projectFolderName = "anio"
)

func GetProjectDataDirectory() (string, error) {
	userDataDir, err := getUserDataDirectory()
	if err != nil {
		return "", fmt.Errorf("failed to get user data directory: %w", err)
	}

	projectDataDir := filepath.Join(userDataDir, projectFolderName)
	if _, err := os.Stat(projectDataDir); err != nil && errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(projectDataDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create project data dir at %s", projectDataDir)
		}
	}

	return projectDataDir, nil
}

func GetProjectConfigDirectory() (string, error) {
	userConfigDir, err := getUserConfigDirectory()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	projectConfigDir := filepath.Join(userConfigDir, projectFolderName)

	if _, err := os.Stat(projectConfigDir); err != nil && errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(projectConfigDir, 0o755); err != nil {
			return "", fmt.Errorf("failed to create project config dir at %s", projectConfigDir)
		}
	}

	return projectConfigDir, nil
}

func getUserDataDirectory() (string, error) {
	switch runtime.GOOS {
	case "linux":
		dataDir := os.Getenv("XDG_DATA_HOME")
		if dataDir != "" {
			return dataDir, nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home dir: %w", err)
		}

		return filepath.Join(homeDir, ".local", "share"), nil

	case "windows":
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config dir: %w", err)
		}

		return cacheDir, nil

	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home dir: %w", err)
		}

		return filepath.Join(homeDir, "Library"), nil

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func getUserConfigDirectory() (string, error) {
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		configFolderPath, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config folder: %w", err)
		}

		return configFolderPath, nil

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
