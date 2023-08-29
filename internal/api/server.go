package api

import (
	"log"
	"net/http"

	controllers "github.com/Masha003/Golang-internship.git/internal/api/controllers/http-transport"
	"github.com/Masha003/Golang-internship.git/internal/api/middleware"
	"github.com/Masha003/Golang-internship.git/internal/config"

	"github.com/gin-gonic/gin"
)

func NewServer(cfg config.Config, userController controllers.UserController) *http.Server {
	log.Print("Creating new server")

	e := gin.Default()
	r := e.Group("/api")

	// Register routes
	registerUserRoutes(r, cfg, userController)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: e,
	}
}

func registerUserRoutes(router *gin.RouterGroup, cfg config.Config, c controllers.UserController) {
	r := router.Group("/users")
	r.POST("/register", c.Register)
	r.POST("/login", c.Login)
	r.GET("/", c.GetAll)
	r.GET("/:id", c.GetById)

	pr := r.Use(middleware.JwtAuth(cfg.Secret))
	pr.GET("/current", c.GetCurrent)
	pr.DELETE("/:id", c.Delete)
	pr.POST("/uploads/:id", c.UploadImage)
	pr.POST("/getimg/:id", c.GetImgByID)
}
