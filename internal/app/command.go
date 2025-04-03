package app

import "fmt"

type Command struct {
	Name      string
	Arguments []string
}

type commands struct {
	Cmds map[string]func(*state, Command) error
}

func (c *commands) register(name string, f func(*state, Command) error) {
	// register a new signature/command
	if c.Cmds == nil {
		c.Cmds = make(map[string]func(*state, Command) error)
	}
	c.Cmds[name] = f
}

func (c *commands) Run(s *state, cmd Command) error {
	_, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return c.Cmds[cmd.Name](s, cmd)
}
