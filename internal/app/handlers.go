package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savisitor15/bootdev-gator/internal/database"
)

func loginHandler(s *state, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username details missing")
	}
	user := cmd.Arguments[0]
	ctx := context.Background()
	_, err := s.db.GetUserByName(ctx, sql.NullString{String: user, Valid: true})
	if err != nil {
		return err
	}
	err = s.appConfig.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Println("login set in config")
	return err
}

func registerHandler(s *state, cmd Command) error {

	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("no username provided")
	}
	name := cmd.Arguments[0]
	uid := uuid.New()
	created_at := time.Now()
	ctx := context.Background()
	params := database.CreateUserParams{
		Name:      sql.NullString{String: name, Valid: true},
		ID:        uid,
		CreatedAt: created_at,
		UpdatedAt: created_at,
	}
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	fmt.Println(user.Name, "created at:", user.CreatedAt)
	return loginHandler(s, Command{Name: "login", Arguments: cmd.Arguments})

}
