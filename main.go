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
	cmds.commandsMap = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("no arguments given")
		os.Exit(1)
	}

	args := os.Args[1:len(os.Args)]
	cmd := command{Name: args[0], Args: args[1:]}
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
