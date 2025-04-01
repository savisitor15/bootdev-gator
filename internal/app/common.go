package app

import (
	config "github.com/savisitor15/bootdev-gator/internal/config"
)

func InitializeState() (state, error) {
	// get the config
	cfg, err := config.Read()
	if err != nil {
		return state{}, err
	}
	return state{appConfig: &cfg}, nil
}

func InitializeCommands() (commands, error) {
	cmds := commands{}
	cmds.register("login", loginHandler)
	return cmds, nil
}
