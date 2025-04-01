package main

import (
	"fmt"
	"os"

	app "github.com/savisitor15/bootdev-gator/internal/app"
)

func buildCommand() (app.Command, error) {
	if len(os.Args[1]) == 0 {
		return app.Command{}, fmt.Errorf("unkown action")
	}
	return app.Command{Name: os.Args[1], Arguments: os.Args[2:len(os.Args)]}, nil
}

func main() {
	state, err := app.InitializeState()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmds, err := app.InitializeCommands()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd, err := buildCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
