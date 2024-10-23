package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ellielle/rssgator/internal/config"
	"github.com/ellielle/rssgator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
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
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	programState := state{
		cfg: &cfg,
		db:  dbQueries,
	}
	coms := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	coms.register("login", handlerLogin)
	coms.register("register", handlerRegister)

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = coms.run(&programState, command{Name: cmdName, Arguments: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
