package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
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
