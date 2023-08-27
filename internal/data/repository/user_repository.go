package repository

import (
	"context"
	"log"
	"time"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Create(t *models.User) error
	Update(t *models.User) error
	Delete(t *models.User) error

	FindByEmail(email string) (models.User, error)
	FindByName(name string) (models.User, error)

	AddToken(ctx context.Context, token string, userID string, expiry time.Duration) error
	GetToken(ctx context.Context, token string) (string, error)
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
	log.Print("Creating new user repository")

	return &userRepository{
		BaseRepository: NewBaseRepository[models.User](db),
		db:             db,
		rdb:            rdb,
	}
}

type userRepository struct {
	BaseRepository[models.User]
	db  *gorm.DB
	rdb *redis.Client
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error

	return user, err
}

func (r *userRepository) FindByName(name string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "name = ?", name).Error

	return user, err
}

func (r *userRepository) AddToken(ctx context.Context, token string, userID string, expiry time.Duration) error {
	return r.rdb.Set(ctx, token, userID, expiry).Err()
}

func (r *userRepository) GetToken(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, token).Result()
}
