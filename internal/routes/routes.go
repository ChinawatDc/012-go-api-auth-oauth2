package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/handlers"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/middlewares"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/repositories"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/services"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/utils"
)

func RegisterRoutes(
	r *gin.Engine,
	db *gorm.DB,
	google *services.GoogleOAuthService,
	jwt *services.JWTService,
	successRedirect string,
) {
	userRepo := repositories.NewUserRepo(db)
	identityRepo := repositories.NewIdentityRepo(db)

	h := handlers.NewOAuthHandler(db, google, userRepo, identityRepo, jwt, successRedirect)

	auth := r.Group("/auth")
	{
		auth.GET("/google/login", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)
	}

	protected := r.Group("/")
	protected.Use(middlewares.AuthRequired(jwt))
	{
		protected.GET("/me", func(c *gin.Context) {
			uidAny, _ := c.Get(middlewares.CtxUserIDKey)
			uid, _ := uidAny.(uint)

			u, err := userRepo.FindByID(uid)
			if err != nil {
				utils.Error(c, 404, "user not found")
				return
			}
			utils.Success(c, gin.H{
				"id":     u.ID,
				"email":  u.Email,
				"name":   u.Name,
				"avatar": u.AvatarURL,
			})
		})
	}
}
