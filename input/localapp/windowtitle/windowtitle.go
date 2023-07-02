package windowtitle

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"anio/config"
	"anio/input/shared"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/rs/zerolog/log"
)

type Poller struct {
	interval       time.Duration
	registeredApps []config.PolledAppConfig
	outputChan     chan<- shared.InputFileInfo
}

func New(interval time.Duration, outputChan chan<- shared.InputFileInfo) *Poller {
	return &Poller{
		interval:   interval,
		outputChan: outputChan,
	}
}

func (poller *Poller) AddApplication(appCfg config.PolledAppConfig) {
	poller.registeredApps = append(poller.registeredApps, appCfg)
}

func (poller *Poller) Start(ctx context.Context) {
	log.Info().Msgf("starting window title poller with an interval of %s", poller.interval)

	ticker := time.NewTicker(poller.interval)

	for {
		select {
		case <-ctx.Done():
			log.Warn().Msgf("closing window title poller due to an abort")
			return
		case <-ticker.C:
			for _, app := range poller.registeredApps {
				var fileName string
				var err error

				switch runtime.GOOS {
				case "windows":
					fileName, err = pollApplicationWindows(app.AppExecutableWindows, app.FilenameMatchRegex)
					if err != nil {
						log.Error().Err(err).Msg("error polling a windows process info")
					}
				case "linux":
					fileName, err = pollApplicationLinux(app.AppExecutableLinux, app.FilenameMatchRegex)
					if err != nil {
						log.Error().Err(err).Msg("error polling a linux process info")
					}
				}

				if fileName == "" {
					continue
				}

				poller.outputChan <- shared.InputFileInfo{FileName: fileName, SourceApplication: app.AppName}
			}
		}
	}
}

func pollApplicationWindows(appName string, matcher *regexp.Regexp) (string, error) {
	filter := fmt.Sprintf("IMAGENAME eq %s", appName)

	cmd := exec.Command("tasklist.exe", "/FI", filter, "/V", "/FO", "csv", "/NH")

	cmdOutput, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read command output for app mpv.exe: %w", err)
	}

	if strings.Contains(string(cmdOutput), "No tasks are running which match") {
		return "", nil
	}

	record, err := csv.NewReader(bytes.NewBuffer(cmdOutput)).Read()
	if err != nil {
		return "", fmt.Errorf("failed to parse command output for app mpv.exe: %w", err)
	}

	if len(record) > 0 {
		fileName := record[len(record)-1]
		cleanedUpFilenameMatches := matcher.FindStringSubmatch(fileName)
		if len(cleanedUpFilenameMatches) > 1 {
			return cleanedUpFilenameMatches[1], nil
		}
		return "", nil
	}

	return "", fmt.Errorf("got an empty record")
}

func pollApplicationLinux(appName string, matcher *regexp.Regexp) (string, error) {
	x, err := xgbutil.NewConn()
	if err != nil {
		return "", fmt.Errorf("can't establish a connection to X server: %w", err)
	}

	windowIDs, err := ewmh.ClientListStackingGet(x)
	if err != nil {
		panic(err)
	}

	for _, windowID := range windowIDs {
		windowName, err := ewmh.WmNameGet(x, windowID)
		if err != nil {
			panic(err)
		}

		if strings.Contains(windowName, appName) {
			cleanedUpFilenameMatches := matcher.FindStringSubmatch(windowName)
			if len(cleanedUpFilenameMatches) > 1 {
				return cleanedUpFilenameMatches[1], nil
			}
		}
	}

	return "", nil
}
