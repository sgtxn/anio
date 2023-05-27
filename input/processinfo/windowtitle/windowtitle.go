package windowtitle

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type InputFileInfo struct {
	FileName          string
	SourceApplication string
}

type LocalApplication struct {
	AppName            string
	AppExecutable      string
	FilenameMatchRegex *regexp.Regexp
}

type WindowTitlePoller struct {
	interval       time.Duration
	registeredApps []LocalApplication
	outputChan     chan<- InputFileInfo
}

func New(interval time.Duration, outputChan chan<- InputFileInfo) *WindowTitlePoller {
	return &WindowTitlePoller{
		interval:   interval,
		outputChan: outputChan,
	}
}

func (poller *WindowTitlePoller) AddApplication(appName, appExecutable, matchRegexp string) {
	poller.registeredApps = append(poller.registeredApps, LocalApplication{
		AppName:            appName,
		AppExecutable:      appExecutable,
		FilenameMatchRegex: regexp.MustCompile(matchRegexp),
	})
}

func (poller *WindowTitlePoller) Start(ctx context.Context) error {
	ticker := time.NewTicker(poller.interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for _, app := range poller.registeredApps {
				fileName, err := pollApplicationWindows(app.AppExecutable, app.FilenameMatchRegex)
				if err != nil {
					log.Error().Err(err).Msg("error polling a windows process info")
				}
				if len(fileName) == 0 {
					continue
				}

				poller.outputChan <- InputFileInfo{FileName: fileName, SourceApplication: app.AppName}
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
		cleanedUpFilename := matcher.FindString(fileName)
		return cleanedUpFilename, nil
	} else {
		return "", fmt.Errorf("got an empty record")
	}
}
