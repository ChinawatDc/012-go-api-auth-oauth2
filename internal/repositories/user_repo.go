package repositories

import (
	"gorm.io/gorm"

	"github.com/ChinawatDc/012-go-api-auth-oauth2/internal/models"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByEmailTx(tx *gorm.DB, email string) (*models.User, error) {
	var u models.User
	if err := tx.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) CreateTx(tx *gorm.DB, u *models.User) error {
	return tx.Create(u).Error
}

func (r *UserRepo) UpdateTx(tx *gorm.DB, u *models.User) error {
	return tx.Save(u).Error
}

func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
