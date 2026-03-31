package service

import (
	"api_orion/model"
	"api_orion/repo"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user model.User) error
	Login(user model.User) (*string, *model.User, error)
	GetProfile(id int) (*model.User, error)
}

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) Register(user model.User) error {
	dbuser, _ := s.userRepo.GetUserByEmail(user.Email)

	if dbuser != nil {
		return errors.New("user already exists")
	}

	return s.userRepo.CreateUser(&user)
}

func (s *userService) Login(user model.User) (token *string, usr *model.User, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, nil, errors.New("wrong email or password")
	}

	expirationTime := time.Now().Add(12 * time.Hour)
	claims := model.Claims{
		UserID: dbUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenJwtString, err := tokenJwt.SignedString(model.JwtKey)
	if err != nil {
		return nil, nil, err
	}

	return &tokenJwtString, dbUser, nil
}

func (s *userService) GetProfile(id int) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
