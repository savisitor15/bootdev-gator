package app

import (
	config "github.com/savisitor15/bootdev-gator/internal/config"
	database "github.com/savisitor15/bootdev-gator/internal/database"
)

type state struct {
	db        *database.Queries
	appConfig *config.Config
}
