package controller

import (
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Masha003/Golang-internship.git/internal/models"
	"github.com/gin-gonic/gin"
)

// http requests, connect with front
type IUserService interface {
	Register(user *models.User) error
	GetUserByID(userID int) (*models.User, error)
	GetUserByName(userName string) (*models.User, error)
	ComparePasswordHash(hash, pass string) error
	Delete(userID int) error
	UploadImage(id int, file *multipart.File, handler *multipart.FileHeader) error
}

type UserController struct {
	userService IUserService
}

func NewUserController(userService IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.User

	// checks error during binding json with user model
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	// checks error during registration
	err := ctrl.userService.Register(&user) //function from user service
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mesage": "user created successfully",
	})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetImgByID(c *gin.Context) {
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.Header("Content-Type", "image/jpeg")
	c.File(user.Image)
}

func (ctrl *UserController) GetUserByName(c *gin.Context) {
	userName := c.Param("name")

	user, err := ctrl.userService.GetUserByName(userName)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Login(c *gin.Context) {
	var loginData struct {
		Name     string `json: "name"`
		Password string `json: "password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	name := loginData.Name
	password := loginData.Password

	user, err := ctrl.userService.GetUserByName(name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized username"})
	}

	err2 := ctrl.userService.ComparePasswordHash(user.Password, password)
	if err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Success!",
	})
}

func (ctrl *UserController) Delete(c *gin.Context) {
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err1 := ctrl.userService.Delete(userID)

	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted from database",
	})

	return
}

func (ctrl *UserController) UploadImage(c *gin.Context) {
	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)
	file, h, err := c.Request.FormFile("path_to_image")

	if err != nil {
		log.Printf("Failed to get image")
	}

	defer file.Close()

	err = ctrl.userService.UploadImage(userID, &file, h)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "image uploaded",
	})
}
