package main

import (
	"anio/config"
	"fmt"

	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("it works!")
	fmt.Println("Loading config.")
	conf, err := config.Load()
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Printf(
		"Hello, %s. This line demonstrates that config "+
			"was successfully loaded.", conf.Name)
}
