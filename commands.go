package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandsMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f := c.commandsMap[cmd.Name]
	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if cmd.Args == nil || len(cmd.Args) == 0 {
		return fmt.Errorf("login command expects username argument")
	}

	username := cmd.Args[0]
	err := s.Cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("Username " + username + " has been saved to configuration")
	return nil
}
