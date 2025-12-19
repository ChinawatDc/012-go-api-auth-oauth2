package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTIssuer     string
	AccessSecret  string
	AccessMinutes int

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	FrontendSuccessRedirect string
}

func Load() Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := Config{
		AppName: viper.GetString("APP_NAME"),
		AppPort: viper.GetString("APP_PORT"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),

		JWTIssuer:    viper.GetString("JWT_ISSUER"),
		AccessSecret: viper.GetString("JWT_ACCESS_SECRET"),

		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),

		FrontendSuccessRedirect: viper.GetString("FRONTEND_SUCCESS_REDIRECT"),
	}

	cfg.AccessMinutes = mustInt(viper.GetString("ACCESS_TOKEN_MINUTES"), 60)

	if cfg.AppPort == "" {
		cfg.AppPort = "8081"
	}
	if cfg.AccessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET is required")
	}
	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" || cfg.GoogleRedirectURL == "" {
		log.Fatal("GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / GOOGLE_REDIRECT_URL are required")
	}
	if cfg.JWTIssuer == "" {
		cfg.JWTIssuer = "012-go-api-auth-oauth2"
	}
	return cfg
}

func mustInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
