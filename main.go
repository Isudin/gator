package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Isudin/gator/internal/config"
	"github.com/Isudin/gator/internal/database"
	_ "github.com/lib/pq"
)

var st state

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st.cfg = &cfg

	db, err := sql.Open("postgres", st.cfg.DbUrl)
	if err != nil {
		fmt.Println("could not connect to the database")
		os.Exit(1)
	}
	dbqueries := database.New(db)
	st.db = dbqueries

	cmds := commands{}
	registerHandlers(&cmds)

	runCommands(&cmds)

	os.Exit(0)
}

func registerHandlers(cmds *commands) {
	cmds.commandsMap = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)
}

func runCommands(cmds *commands) {
	if len(os.Args) < 2 {
		fmt.Println("no arguments given")
		os.Exit(1)
	}

	args := os.Args[1:len(os.Args)]
	commandName := args[0]
	if cmds.commandsMap[commandName] == nil {
		fmt.Printf("command '%v' not found\n", commandName)
		os.Exit(1)
	}

	cmd := command{Name: commandName, Args: args[1:]}
	err := cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
