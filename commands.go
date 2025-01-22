package main

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
