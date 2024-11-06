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
	coms.register("reset", handlerReset)
	coms.register("users", handlerGetUsers)
	coms.register("agg", handlerAggregate)
	coms.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	coms.register("feeds", handlerGetFeeds)
	coms.register("follow", middlewareLoggedIn(handlerAddFollow))
	coms.register("following", middlewareLoggedIn(handlerFollowing))
	coms.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = coms.run(&programState, command{Name: cmdName, Arguments: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}

// register adds commands to the commands struct
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run ensures a command has been registered, and then runs it
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		log.Fatalf("%s is not a command", cmd.Name)
	}
	return f(s, cmd)
}
