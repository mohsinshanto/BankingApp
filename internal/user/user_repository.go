package user

import (
	"banking/models"
	"errors"

	appError "banking/errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
func (r *userRepository) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email=?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appError.ErrInvalidCredentials
		}
		return nil, err

	}
	return &user, nil
}
