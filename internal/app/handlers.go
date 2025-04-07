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

func usersHandler(s *state, _ Command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Users:")
	fmt.Println("----------")
	for _, elm := range users {
		if elm.Name.String == s.appConfig.CurrentUserName {
			elm.Name.String = elm.Name.String + " (current)"
		}
		fmt.Println(elm.Name.String)
	}
	return nil
}

func resetHandler(s *state, _ Command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		return err
	}
	return nil
}

func aggHandler(s *state, _ Command) error {
	var ctx context.Context = context.Background()
	res, err := fetchFeed(ctx, s.appConfig.FeedUrl)
	if err != nil {
		return err
	}
	cleanFeedOutput(res)
	fmt.Println("channel:", res.Channel.Title)
	fmt.Println(res.Channel.Description)
	fmt.Println("items:")
	for _, itm := range res.Channel.Item {
		fmt.Println("title:", itm.Title)
		fmt.Println(itm.Description)
	}

	return nil
}

func addfeedHandler(s *state, cmd Command) error {
	var ctx context.Context = context.Background()
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("not enough arguments")
	}
	ts := time.Now()
	name := cmd.Arguments[0]
	url := cmd.Arguments[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      sql.NullString{String: name, Valid: true},
		Url:       sql.NullString{String: url, Valid: true},
		UserID:    uuid.NullUUID{UUID: s.currentUser.ID, Valid: true},
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}
	return followHandler(s, Command{Name: "follow", Arguments: []string{feed.Url.String}})
}

func feedsHandler(s *state, _ Command) error {
	ctx := context.Background()
	output, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Feed | URL | UserName")
	for _, elm := range output {
		fmt.Println(elm.Name.String, "|", elm.Url.String, "|", elm.Username.String)
	}
	return nil
}

func followHandler(s *state, cmd Command) error {
	ctx := context.Background()

	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("missing url for follow")
	}
	url := cmd.Arguments[0]
	feed, err := s.db.GetFeedsByUrl(ctx, sql.NullString{String: url, Valid: true})
	if err != nil {
		return err
	}
	ts := time.Now()
	params := database.CreateFeedFollowParams{
		CreatedAt: ts,
		UpdatedAt: ts,
		UserID:    uuid.NullUUID{UUID: s.currentUser.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	}
	follow, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Println("new follow registered for", s.appConfig.CurrentUserName, "on feed", follow.FeedName.String)
	return nil
}

func followingHandler(s *state, _ Command) error {
	ctx := context.Background()
	follows, err := s.db.GetFeedFollowsForUser(ctx, s.currentUser.ID)
	if err != nil {
		return err
	}
	fmt.Println("Follows for user:", s.currentUser.Name.String)
	fmt.Println("User has", len(follows), "subscribed feeds")
	fmt.Println("ID", "|", "feed", "|", "created at")
	for _, elm := range follows {
		fmt.Println(elm.ID, "|", elm.FeedName, "|", elm.CreatedAt)
	}
	return nil
}
