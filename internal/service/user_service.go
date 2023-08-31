package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/Masha003/Golang-internship.git/internal/auth"
	"github.com/Masha003/Golang-internship.git/internal/config"
	"github.com/Masha003/Golang-internship.git/internal/data/repository"
	"github.com/Masha003/Golang-internship.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user models.RegisterUser) (models.Token, error)
	Login(user models.LoginUser) (models.Token, error)
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	FindByName(userName string) (models.User, error)
	Delete(id string) error
	UploadImage(id string, file *multipart.File, handler *multipart.FileHeader) error
}

func NewUserService(userRepo repository.UserRepository, cfg config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type userService struct {
	userRepo repository.UserRepository
	cfg      config.Config
}

func (s *userService) Register(registerUser models.RegisterUser) (models.Token, error) {
	var token models.Token

	user, err := s.userRepo.FindByEmail(registerUser.Email)
	if err == nil {
		return token, errors.New("user with this email already exists")
	}

	pass, err := generatePasswordHash(registerUser.Password)
	if err != nil {
		return token, err
	}

	user.Email = registerUser.Email
	user.Name = registerUser.Name
	user.Password = pass

	err = s.userRepo.Create(&user)
	if err != nil {
		return token, err
	}

	tokenString, refreshTokenString, err := auth.GenerateJWT(user.ID, s.cfg.TokenLifespan, s.cfg.Secret)
	if err != nil {
		return token, err
	}

	token.Token = tokenString
	token.RefreshToken = refreshTokenString
	token.User = user

	return token, nil
}

func (s *userService) Login(loginUser models.LoginUser) (models.Token, error) {
	var token models.Token

	user, err := s.userRepo.FindByEmail(loginUser.Email)
	if err != nil {
		return token, err
	}

	err = comparePasswordHash(user.Password, loginUser.Password)
	if err != nil {
		return token, err
	}

	tokenString, refreshTokenString, err := auth.GenerateJWT(user.ID, s.cfg.TokenLifespan, s.cfg.Secret)
	if err != nil {
		return token, err
	}

	token.Token = tokenString
	token.RefreshToken = refreshTokenString
	token.User = user

	return token, nil
}

func generatePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}

	return string(hash), nil
}

func comparePasswordHash(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

func (s *userService) FindAll(query models.PaginationQuery) ([]models.User, error) {
	return s.userRepo.FindAll(query)
}

func (s *userService) FindById(id string) (models.User, error) {
	return s.userRepo.FindById(id)
}

func (s *userService) FindByName(userName string) (models.User, error) {
	return s.userRepo.FindByName(userName)
}

func (s *userService) Delete(id string) error {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(&user)
}

func (s *userService) UploadImage(userID string, file *multipart.File, handler *multipart.FileHeader) error {
	filePath := "./uploads/" + handler.Filename
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("unable to upload image 1")
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, *file)
	if err != nil {
		log.Printf("unable to upload image 2")
		return err
	}

	user, err := s.userRepo.FindById(userID)
	if err != nil {
		log.Printf("unable to upload image 3: %v\n", err)
		return err
	}

	user.Image = filePath

	return s.userRepo.Update(&user)
}
