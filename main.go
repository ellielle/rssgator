package main

import (
	"fmt"

	"github.com/ellielle/rssgator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	cfg.SetUser()

	cfg, err = config.Read()
	fmt.Printf("%v", cfg)
}

func handlerLogin(st *state, cmd command) error {

	return nil
}
