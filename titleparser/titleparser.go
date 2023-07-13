package titleparser

import (
	"fmt"
	"strconv"

	"anio/input/shared"

	"github.com/nssteinbrenner/anitogo"
	"github.com/rs/zerolog/log"
)

type Titleparser struct {
	input  <-chan shared.PlaybackFileInfo
	output chan<- shared.PlaybackAnimeDetails
}

func New(input <-chan shared.PlaybackFileInfo, output chan<- shared.PlaybackAnimeDetails) *Titleparser {
	return &Titleparser{
		input:  input,
		output: output,
	}
}

func (tp *Titleparser) Start() {
	go tp.runDataLoop()
}

func (tp *Titleparser) runDataLoop() {
	for inputfile := range tp.input {
		log.Debug().Any("titleParserInput", inputfile).Msg("got input")
		elements := anitogo.Parse(inputfile.FileName, anitogo.DefaultOptions)

		if err := validateElements(elements); err != nil {
			log.Error().Err(err).Str("file", inputfile.FileName).Msg("error while parsing filename")
			continue
		}

		result, err := buildPlaybackDetails(elements)
		if err != nil {
			log.Error().Err(err).Str("file", inputfile.FileName).Msg("error while building titleparser output")
			continue
		}

		tp.output <- result
	}
}

func validateElements(elements *anitogo.Elements) error {
	if elements.AnimeTitle == "" {
		return fmt.Errorf("failed to extract title")
	}

	if len(elements.EpisodeNumber) == 0 {
		return fmt.Errorf("failed to extract episode number")
	}

	return nil
}

func buildPlaybackDetails(elements *anitogo.Elements) (shared.PlaybackAnimeDetails, error) {
	var result shared.PlaybackAnimeDetails
	result.Title = elements.AnimeTitle

	epNumbers := make([]int, 0, len(elements.EpisodeNumber))
	for _, epNoString := range elements.EpisodeNumber {
		epNo, err := strconv.ParseInt(epNoString, 10, 64)
		if err != nil {
			return result, fmt.Errorf("couldn't parse episode number %s as a number", epNoString)
		}

		epNumbers = append(epNumbers, int(epNo))
	}

	// only return the last one, cause if we got multiple episodes in one file - we're only interested in the last one
	result.Episode = epNumbers[len(epNumbers)-1]

	return result, nil
}
