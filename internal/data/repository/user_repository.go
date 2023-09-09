package repository

import (
	"context"
	"log"
	"time"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewUserMongoRepository(rdb *redis.Client, mongodb *mongo.Database) UserRepository {
	log.Print("Creating new user mongo repository")

	return &userMongoRepository{
		rdb:        rdb,
		mongodb:    mongodb,
		collection: mongodb.Collection("users"),
	}
}

type userMongoRepository struct {
	rdb        *redis.Client
	mongodb    *mongo.Database
	collection *mongo.Collection
}

func (r *userMongoRepository) FindAll(query models.PaginationQuery) ([]models.User, error) {
	var users []models.User

	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return users, err
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *userMongoRepository) FindById(id string) (models.User, error) {
	var user models.User

	return user, nil
}

func (r *userMongoRepository) Create(t *models.User) error {
	return nil
}

func (r *userMongoRepository) Update(t *models.User) error {
	return nil
}

func (r *userMongoRepository) Delete(t *models.User) error {
	return nil
}

func (r *userMongoRepository) FindByEmail(email string) (models.User, error) {
	var user models.User

	return user, nil
}

func (r *userMongoRepository) FindByName(name string) (models.User, error) {
	var user models.User

	return user, nil
}

func (r *userMongoRepository) AddToken(ctx context.Context, token string, userID string, expiry time.Duration) error {
	return r.rdb.Set(ctx, token, userID, expiry).Err()
}

func (r *userMongoRepository) GetToken(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, token).Result()
}
