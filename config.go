package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Port string `default:"8080"`
}

func newConfig() (*config, error) {
	c := &config{}
	err := envconfig.Process("", c)
	return c, err
}
