package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConfigFilename = "config.json" // can't use other names
var supportedOS = []string{"windows", "linux", "darwin"}

var (
	workDirectory, _ = os.Getwd()
	testDirectory    = filepath.Join(workDirectory, "tests")
	configFilePath   = filepath.Join(testDirectory, testConfigFilename)
)

// [TestExists]: Just testing if the helper function works properly.
func TestExists(t *testing.T) {
	os.Mkdir(testDirectory, 0o644)
	assert.True(t, exists(testDirectory))
	os.Remove(testDirectory)
	assert.False(t, exists(testDirectory))
}

// [testConfigContent]: Helper Function that checks the conf.
func testConfigContent(t *testing.T, conf *Config) {
	assert.NotNil(t, conf)
	assert.NotEmpty(t, conf.Name)
	assert.True(t, exists(configFilePath)) // uses a global var
}

// TestCase: Everything all right.
func TestCase_EverythingGood(t *testing.T) {
	// check for test config existence
	if exists(configFilePath) {
		os.Remove(configFilePath)
	}

	// first test, creating the new cfg and storing in file
	cfg_new, err := Load(configFilePath)
	assert.NoError(t, err)
	testConfigContent(t, cfg_new)

	// second test, loading the cfg from file
	cfg_loaded, err := Load(configFilePath)
	assert.NoError(t, err)
	testConfigContent(t, cfg_loaded)

	os.RemoveAll(testDirectory)
}

// TestCase: Function is provided with an inaccesible path.
func TestCase_WrongFolder_ExpectedError(t *testing.T) {
	wrongDirectory := "A://bcdefg/"
	wrongFilePath := filepath.Join(wrongDirectory, testConfigFilename)
	_, err := Load(wrongFilePath)
	assert.Error(t, err)
}

// TestCase: Function tries to load cfg from a file that is not Json at all.
func TestCase_UnreadableJson_ExpectedError(t *testing.T) {
	data := []byte("test")
	_ = os.Mkdir(testDirectory, 0o755)
	file, err := os.Create(configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	file.Write(data)
	file.Close()

	_, err = Load(configFilePath)
	assert.Error(t, err)

	os.RemoveAll(testDirectory)
}
