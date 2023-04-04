package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string // develop, staging, production
	
	// context timeout in seconds
	CtxTimeout int

	LogLevel string
	HTTPPort string
	
	UserServiceHost string
	UserServicePort string
	PostServiceHost  string
	PostServicePort  string
	CommentServiceHost string
	CommentServicePort string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "localhost"))
	c.UserServicePort = cast.ToString(getOrReturnDefault("USER_SERVICE_PORT", "8000"))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "localhost"))
	c.PostServicePort = cast.ToString(getOrReturnDefault("POST_SERVICE_PORT", "8010"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "localhost"))
	c.CommentServicePort = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_PORT", "8020"))


	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
