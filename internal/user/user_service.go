package user

import (
	"banking/config"
	appError "banking/errors"
	"banking/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(input *RegistrationInput) error
	Login(input *LoginInput) (string, error)
}
type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (s *userService) Register(input *RegistrationInput) error {
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
func (s *userService) Login(input *LoginInput) (string, error) {
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
