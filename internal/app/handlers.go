package app

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

func aggHandlerHelper(s *state) error {
	var ctx context.Context = context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	res, err := fetchFeed(ctx, nextFeed.Url.String)
	if err != nil {
		return err
	}
	ts := time.Now()
	s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{ID: nextFeed.ID, UpdatedAt: ts})
	cleanFeedOutput(res)
	fmt.Println("channel:", res.Channel.Title)
	fmt.Println(res.Channel.Description)
	fmt.Println("items:")
	for _, itm := range res.Channel.Item {
		ts = time.Now()
		pubt, err := time.Parse(time.RFC1123, itm.PubDate)
		if err != nil {
			fmt.Println(err)
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   ts,
			Title:       sql.NullString{String: itm.Title, Valid: true},
			Description: sql.NullString{String: itm.Description, Valid: true},
			Url:         sql.NullString{String: itm.Link, Valid: true},
			PublishedAt: sql.NullTime{Time: pubt, Valid: true},
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func aggHandler(s *state, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("too few arguments, require interval: 1s, 1m, 1h")
	}
	time_between_reqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every", time_between_reqs)
	// ticker := time.NewTicker(time_between_reqs)
	// for ; ; <-ticker.C {
	// 	aggHandlerHelper(s)
	// }
	return aggHandlerHelper(s)
}

func addfeedHandler(s *state, cmd Command, user database.User) error {
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
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	feed, err := s.db.CreateFeed(ctx, params)
	if err != nil {
		return err
	}
	return followHandler(s, Command{Name: "follow", Arguments: []string{feed.Url.String}}, user)
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

func followHandler(s *state, cmd Command, user database.User) error {
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
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	}
	follow, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Println("new follow registered for", s.appConfig.CurrentUserName, "on feed", follow.FeedName.String)
	return nil
}

func followingHandler(s *state, _ Command, user database.User) error {
	ctx := context.Background()
	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Println("Follows for user:", user.Name.String)
	fmt.Println("User has", len(follows), "subscribed feeds")
	fmt.Println("ID", "|", "feed", "|", "created at")
	for _, elm := range follows {
		fmt.Println(elm.ID, "|", elm.FeedName, "|", elm.CreatedAt)
	}
	return nil
}

func unfollowHandler(s *state, cmd Command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("url to unfollow missing")
	}
	url := cmd.Arguments[0]
	err := s.db.DeleteFeedFollowsForUser(ctx,
		database.DeleteFeedFollowsForUserParams{Name: user.Name, Url: sql.NullString{String: url, Valid: true}})
	return err
}

func browseHandler(s *state, cmd Command, user database.User) error {
	var limit int = 0
	var err error
	ctx := context.Background()
	if len(cmd.Arguments) >= 1 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return err
		}
	}else{
		limit = 2
	}
	res, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit: int32(limit),
	})
	if err != nil{
		return err
	}
	for _, elm := range res{
		fmt.Println("Published at", elm.PublishedAt.Time)
		fmt.Println("Title", elm.Title.String)
		fmt.Println("Link", elm.Url.String)
		fmt.Println("Description")
		fmt.Println(elm.Description.String)		
	}
	return nil
}
