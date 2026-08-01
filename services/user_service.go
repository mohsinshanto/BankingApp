package services

import (
	"banking/config"
	"banking/dto"
	appError "banking/errors"
	"banking/models"
	"banking/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(input *dto.RegistrationInput) error
	Login(input *dto.LoginInput) (string, error)
}
type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (s *userService) Register(input *dto.RegistrationInput) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dbUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashPassword),
	}
	if err := s.repo.Create(&dbUser); err != nil {
		return err
	}

	return nil

}
func (s *userService) Login(input *dto.LoginInput) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", appError.ErrInvalidCredentials
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
