package service

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// business logic

type IUserRepository interface {
	Insert(user *models.User) error
	CheckIfEmailExist(mail string) bool
	GetUserByID(id int) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	Delete(id int) error
	UpdateUserImage(id int, path string) error
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserSerice(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Register(user *models.User) error {
	if svc.userRepo.CheckIfEmailExist(user.Email) {
		return errors.New("User already exists")
	}

	hash, err := svc.GeneratePasswordHash(user.Password)

	if err != nil {
		return errors.New("Can't register user")
	}

	user.Password = hash

	return svc.userRepo.Insert(user)

}

func (svc *UserService) GeneratePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return "", err
	}

	return string(hash), nil
}

func (svc *UserService) ComparePasswordHash(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	return nil
}

func (svc *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := svc.userRepo.GetUserByID(id)

	if err != nil {
		log.Printf("failed to get user by id: %v\n", err)
		return nil, err
	}

	return user, nil
}

func (svc *UserService) GetUserByName(userName string) (*models.User, error) {
	user, err := svc.userRepo.GetUserByName(userName)

	if err != nil {
		log.Printf("Failed to find user by name: %v\n", err)
		return nil, err
	}

	return user, nil
}

func (svc *UserService) Delete(userID int) error {
	err := svc.userRepo.Delete(userID)

	if err != nil {
		log.Printf("Failed to delete user from database: %v\n", err)
		return err
	}

	return nil
}

func (svc *UserService) UploadImage(userID int, file *multipart.File, handler *multipart.FileHeader) error {
	filePath := "./uploads/" + handler.Filename
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("unable to upload image")
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, *file)
	if err != nil {
		log.Printf("unable to upload image")
		return err
	}
	err = svc.userRepo.UpdateUserImage(userID, filePath)

	if err != nil {
		return err
	}

	return nil
}
