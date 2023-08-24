package main

import (
	"log"

	"github.com/Masha003/Golang-internship.git/internal/controller"
	"github.com/Masha003/Golang-internship.git/internal/models"
	"github.com/Masha003/Golang-internship.git/internal/repository"
	"github.com/Masha003/Golang-internship.git/internal/service"
	"github.com/Masha003/Golang-internship.git/pkg/database/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserSerice(userRepo)
	userController := controller.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.POST("/", userController.Register)
		userRouter.GET("/:id", userController.GetUserByID)
		userRouter.GET("/login", userController.Login)
		userRouter.DELETE("/:id", userController.Delete)
		userRouter.POST("/img/:id", userController.UploadImage)
		userRouter.POST("/get/:id", userController.GetImgByID)

		_ = router.Run(":8888")
	}
}

func init() {
	db = postgres.NewDBConnection()
	err := db.AutoMigrate(models.User{})

	if err != nil {
		log.Fatalf("Failed to migrate user model \n")
	}

}
