package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/models"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/repositories"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/services"
	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/utils"
)

type OAuthHandler struct {
	db              *gorm.DB
	google          *services.GoogleOAuthService
	users           *repositories.UserRepo
	idents          *repositories.IdentityRepo
	jwt             *services.JWTService
	successRedirect string
}

func NewOAuthHandler(
	db *gorm.DB,
	google *services.GoogleOAuthService,
	users *repositories.UserRepo,
	idents *repositories.IdentityRepo,
	jwt *services.JWTService,
	successRedirect string,
) *OAuthHandler {
	return &OAuthHandler{
		db:              db,
		google:          google,
		users:           users,
		idents:          idents,
		jwt:             jwt,
		successRedirect: successRedirect,
	}
}

func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state := randomState(32)

	// demo (ง่าย): เก็บ state ใน cookie 5 นาที
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	url := h.google.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	cookieState, _ := c.Cookie("oauth_state")
	if cookieState == "" || state == "" || cookieState != state {
		utils.Error(c, http.StatusBadRequest, "invalid oauth state")
		return
	}
	if code == "" {
		utils.Error(c, http.StatusBadRequest, "missing code")
		return
	}

	token, err := h.google.Exchange(c.Request.Context(), code)
	if err != nil {
		utils.ErrorWithDetail(c, http.StatusUnauthorized, "exchange code failed", err.Error())
		return
	}

	info, err := h.google.FetchUserInfo(c.Request.Context(), token)
	if err != nil {
		utils.ErrorWithDetail(c, http.StatusUnauthorized, "fetch userinfo failed", err.Error())
		return
	}

	if info.Email == "" || info.Sub == "" {
		utils.Error(c, http.StatusUnauthorized, "invalid google userinfo")
		return
	}

	var user *models.User

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 1) หา user ตาม email
		u, findErr := h.users.FindByEmailTx(tx, info.Email)
		if findErr == nil {
			u.Name = info.Name
			u.AvatarURL = info.Picture
			if err := h.users.UpdateTx(tx, u); err != nil {
				return err
			}
			user = u
			return h.idents.UpsertTx(tx, u.ID, "google", info.Sub, info.Email)
		}

		// 2) ไม่เจอ → สร้างใหม่
		newUser := &models.User{
			Email:     info.Email,
			Name:      info.Name,
			AvatarURL: info.Picture,
		}
		if err := h.users.CreateTx(tx, newUser); err != nil {
			return err
		}
		user = newUser
		return h.idents.UpsertTx(tx, newUser.ID, "google", info.Sub, info.Email)
	})
	if err != nil {
		utils.ErrorWithDetail(c, http.StatusInternalServerError, "upsert user failed", err.Error())
		return
	}

	access, err := h.jwt.NewAccessToken(user.ID)
	if err != nil {
		utils.ErrorWithDetail(c, http.StatusInternalServerError, "create jwt failed", err.Error())
		return
	}

	// Option A: ส่ง JSON
	// utils.Success(c, gin.H{"access_token": access, "token_type": "Bearer"})
	// return

	// Option B: redirect ไป frontend
	if h.successRedirect != "" {
		c.Redirect(http.StatusFound, h.successRedirect+"?access_token="+access)
		return
	}

	utils.Success(c, gin.H{"access_token": access, "token_type": "Bearer"})
}

func randomState(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
