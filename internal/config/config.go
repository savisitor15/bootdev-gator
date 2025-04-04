package config

type Config struct {
	DbURL           string `json:"db_url"`
	FeedUrl			string `json:"feed_url"`
	CurrentUserName string `json:"current_user_name"`
}
