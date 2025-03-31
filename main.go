package main

import (
	config "github.com/savisitor15/bootdev-gator/internal/config"
	"fmt"
)

func main(){
	var cfg config.Config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	cfg.SetUser("lane")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}