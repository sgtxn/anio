package main

import (
	"anio/config"
	"fmt"
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
			"was succesfully loaded.", conf.Name)
}
