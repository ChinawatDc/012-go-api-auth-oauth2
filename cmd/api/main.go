package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/config"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/db"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/routes"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/services"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal("db connect error:", err)
	}
	log.Println("GOOGLE_CLIENT_ID:", cfg.GoogleClientID)
	log.Println("GOOGLE_REDIRECT_URL:", cfg.GoogleRedirectURL)
	log.Println("GOOGLE_CLIENT_SECRET len:", len(cfg.GoogleClientSecret))
	log.Println("GOOGLE_CLIENT_SECRET:", cfg.GoogleClientSecret)

	googleSvc := services.NewGoogleOAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	jwtSvc := services.NewJWTService(cfg)

	r := gin.Default()
	routes.RegisterRoutes(r, database, googleSvc, jwtSvc, cfg.FrontendSuccessRedirect)

	addr := ":" + cfg.AppPort
	log.Println("server running at", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
