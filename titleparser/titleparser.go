package titleparser

import (
	"fmt"

	"anio/input/shared"

	"github.com/nssteinbrenner/anitogo"
	"github.com/rs/zerolog/log"
)

type Titleparser struct {
	input <-chan shared.InputFileInfo
}

func New(input <-chan shared.InputFileInfo) *Titleparser {
	return &Titleparser{
		input: input,
	}
}

func (tp *Titleparser) Start() {
	for inputfile := range tp.input {
		elements := anitogo.Parse(inputfile.FileName, anitogo.DefaultOptions)

		if err := validateElements(elements); err != nil {
			log.Error().Err(err).Str("file", inputfile.FileName).Msg("error while parsing filename")
			continue
		}
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
