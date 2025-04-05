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
	user, err := s.db.GetUserByName(ctx, sql.NullString{String: "_g_invalid", Valid: true})
	if err != nil {
		params := database.CreateUserParams{
			Name:      sql.NullString{String: "_g_invalid", Valid: true},
			ID:        uuid.New(),
			CreatedAt: ts,
			UpdatedAt: ts,
		}
		s.db.CreateUser(ctx, params)
	}
	_, err = s.db.GetFeedByUser(ctx, uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		params := database.CreateFeedParams{
			Name:      sql.NullString{String: "_g_invalid", Valid: true},
			ID:        uuid.New(),
			CreatedAt: ts,
			UpdatedAt: ts,
		}
		s.db.CreateFeed(ctx, params)
	}
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
	ctx := context.Background()
	user, _ := dbq.GetUserByName(ctx, sql.NullString{String: cfg.CurrentUserName, Valid: true})
	return state{
		appConfig:   &cfg,
		db:          dbq,
		currentUser: &user,
	}, nil
}

func initializeCommands() (commands, error) {
	cmds := commands{}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", addfeedHandler)
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
