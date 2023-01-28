package main

import (
	"fmt"

	"anio/config"
)

func main() {
	fmt.Println("it works!")
	fmt.Println("Loading config.")
	conf, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(
		"Hello, %s. This line demonstrates that config "+
			"was successfully loaded.", conf.Name)
}
