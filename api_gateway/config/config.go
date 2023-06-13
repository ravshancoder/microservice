package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string // develop, staging, production

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	// context timeout in seconds
	CtxTimeout int

	SiginKey string

	LogLevel string
	HTTPPort string

	UserServiceHost string
	UserServicePort string

	AuthConfigPath string

	PostServiceHost string
	PostServicePort string

	CommentServiceHost string
	CommentServicePort string

	RedisHost string
	RedisPort string
}

func Load() Config {
	c := Config{}

	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "ravshan"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "r"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "database"))
	c.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5434"))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "userdb"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.SiginKey = cast.ToString(getOrReturnDefault("SIGNING_KEY", "ravshanSignIn"))

	c.AuthConfigPath = cast.ToString(getOrReturnDefault("CASBIN_CONFIG_PATH", "./config/rback_model.conf"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8030"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToString(getOrReturnDefault("USER_SERVICE_PORT", "8000"))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
	c.PostServicePort = cast.ToString(getOrReturnDefault("POST_SERVICE_PORT", "8010"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment_service"))
	c.CommentServicePort = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_PORT", "8020"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redisdb"))
	c.RedisPort = cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
