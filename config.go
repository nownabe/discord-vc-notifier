package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	DiscordChannelID string `split_words:"true"`
	Port             string `default:"8080"`
	SlackChannel     string `split_words:"true"`
	SlackIconEmoji   string `split_words:"true"`
	SlackUsername    string `split_words:"true"`
	SlackWebhookURL  string `split_words:"true"`
}

func newConfig() (*config, error) {
	c := &config{}
	err := envconfig.Process("", c)
	return c, err
}
