package repositories

import (
	"gorm.io/gorm"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/models"
)

type IdentityRepo struct{ db *gorm.DB }

func NewIdentityRepo(db *gorm.DB) *IdentityRepo { return &IdentityRepo{db: db} }

func (r *IdentityRepo) UpsertTx(tx *gorm.DB, userID uint, provider, providerUID, email string) error {
	var ident models.OAuthIdentity
	err := tx.Where("provider = ? AND provider_uid = ?", provider, providerUID).First(&ident).Error
	if err == nil {
		ident.UserID = userID
		ident.Email = email
		return tx.Save(&ident).Error
	}

	ident = models.OAuthIdentity{
		UserID:      userID,
		Provider:    provider,
		ProviderUID: providerUID,
		Email:       email,
	}
	return tx.Create(&ident).Error
}
