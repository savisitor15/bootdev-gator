package config

import (
	"fmt"
)

func (c Config) SetUser (user string) error {
	if len(user) == 0 {
		return fmt.Errorf("user must be defined")
	}
	c.CurrentUserName = user
	err := write(c)
	return err
}
