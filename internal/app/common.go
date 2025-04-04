package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	config "github.com/savisitor15/bootdev-gator/internal/config"
	database "github.com/savisitor15/bootdev-gator/internal/database"
)

func initDefaultUser(s *state) {
	ts := time.Now()
	ctx := context.Background()
	params := database.CreateUserParams{
		Name:      sql.NullString{String: "_g_invalid", Valid: true},
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	s.db.CreateUser(ctx, params)
}

func initializeDatabase(cfg *config.Config) (*database.Queries, error) {
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		return nil, err
	}
	return database.New(db), nil
}

func initializeState() (state, error) {
	// get the config
	cfg, err := config.Read()
	if err != nil {
		return state{}, err
	}
	dbq, err := initializeDatabase(&cfg)
	if err != nil {
		return state{}, err
	}
	return state{
		appConfig: &cfg,
		db:        dbq,
	}, nil
}

func initializeCommands() (commands, error) {
	cmds := commands{}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	return cmds, nil
}

func InitializeApp() (state, commands, error) {
	st, err := initializeState()
	if err != nil {
		return state{}, commands{}, fmt.Errorf("unable to initialize state: %w", err)
	}
	// blind fire this
	initDefaultUser(&st)
	cmds, err := initializeCommands()
	if err != nil {
		fmt.Println(err)
		return state{}, commands{}, fmt.Errorf("unable to initialize command structure: %w", err)
	}
	return st, cmds, nil
}
