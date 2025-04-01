package app

import (
	"fmt"
)

func loginHandler(s *state, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username details missing")
	}
	user := cmd.Arguments[0]
	err := s.appConfig.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Println("login set in config")
	return err
}
