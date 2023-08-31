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
	UploadImage(ctx *gin.Context)
	GetImgByID(ctx *gin.Context)
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

func (c *userController) UploadImage(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	file, h, err := ctx.Request.FormFile("path_to_image")
	if err != nil {
		log.Printf("Failed to get image")
	}

	defer file.Close()

	err = c.service.UploadImage(userIDStr, &file, h)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "image uploaded",
	})
}

func (c *userController) GetImgByID(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := c.service.FindById(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Header("Content-Type", "image/jpeg")
	ctx.File(user.Image)
}
