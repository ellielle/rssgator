package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ellielle/rssgator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	programState := state{
		cfg: &cfg,
	}
	coms := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	coms.register("login", handlerLogin)

	//
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = coms.run(&programState, command{Name: cmdName, Arguments: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
