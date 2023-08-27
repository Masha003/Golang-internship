package controllers

import (
	"log"
	"net/http"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"github.com/Masha003/Golang-internship.git/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetCurrent(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Delete(ctx *gin.Context)
	// UploadImage(ctx *gin.Context)
}

func NewUserController(service service.UserService) UserController {
	log.Print("Creating new user controller")

	return &userController{
		service: service,
	}
}

type userController struct {
	service service.UserService
}

func (c *userController) GetAll(ctx *gin.Context) {
	query := models.PaginationQuery{}
	err := ctx.BindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := c.service.FindAll(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *userController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) GetCurrent(ctx *gin.Context) {
	id := ctx.GetString("user_id")

	user, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Register(ctx *gin.Context) {
	var user models.RegisterUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Register(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, token)
}

func (c *userController) Login(ctx *gin.Context) {
	var user models.LoginUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// func (ctrl *UserController) GetImgByID(c *gin.Context) {
// 	userIDStr := c.Param("id")

// 	userID, err := strconv.Atoi(userIDStr)

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error":   true,
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	user, err := ctrl.userService.GetUserByID(userID)

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
// 			"error":   true,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	c.Header("Content-Type", "image/jpeg")
// 	c.File(user.Image)
// }

// func (ctrl *UserController) UploadImage(c *gin.Context) {
// 	userIDStr := c.Param("id")

// 	userID, err := strconv.Atoi(userIDStr)
// 	file, h, err := c.Request.FormFile("path_to_image")

// 	if err != nil {
// 		log.Printf("Failed to get image")
// 	}

// 	defer file.Close()

// 	err = ctrl.userService.UploadImage(userID, &file, h)

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error":   true,
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "image uploaded",
// 	})
// }
