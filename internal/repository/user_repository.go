package repository

import (
	"errors"
	"log"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"gorm.io/gorm"
)

//  connect with database

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Insert(user *models.User) error {
	err := repo.dbClient.Debug().Model(models.User{}).Create(user).Error

	if err != nil {
		log.Printf("Failed to insert user in database: %v\n", err)
		return err
	}
	return nil
}

func (repo *UserRepository) CheckIfEmailExist(mail string) bool {
	var user models.User
	err := repo.dbClient.Debug().Model(models.User{}).Find(&user).Where("email = ?", mail).Error
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (repo *UserRepository) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := repo.dbClient.Debug().Model(models.User{}).First(user, userID).Error

	if err != nil {
		log.Printf("Failed to get user by ID %v\n", err)
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) GetUserByName(userName string) (*models.User, error) {
	user := &models.User{}

	err := repo.dbClient.Debug().Model(models.User{}).Where("name = ?", userName).Error

	if err != nil {
		log.Printf("Failed to find user with this name: %v\n", err)
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) Delete(userID int) error {
	err := repo.dbClient.Debug().Model(models.User{}).Where("id = ?", userID).Delete(&models.User{}).Error

	if err != nil {
		log.Printf("Unable to delete user: %v\n", err)
		return err
	}
	return nil
}

func (repo *UserRepository) UpdateUserImage(userId int, imagePath string) error {
	err := repo.dbClient.Debug().Model(models.User{}).Where("id = ?", userId).Update("image", imagePath).Error

	if err != nil {
		log.Printf("unable to see image: %v\n", err)
		return err
	}
	return nil
}
