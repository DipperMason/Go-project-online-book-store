package config

import "config"

type Config struct {
	HTTPAddr       string
	JWTSecret      string
	AccessTokenTTL int
}

func NewConfig() Config {
	err := config.LoadDotEnv(".env")
	if err != nil {
		panic(err)
	}

	return Config{
		HTTPAddr:       config.Get("HTTP_ADDR", ":8001"),
		JWTSecret:      config.Get("JWT_SECRET", "hihihaha"),
		AccessTokenTTL: config.GetInt("ACCESS_TOKEN_TTL", 86400),
	}
}
