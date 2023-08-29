package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Masha003/Golang-internship.git/internal/api"
	"github.com/Masha003/Golang-internship.git/internal/api/controllers"
	"github.com/Masha003/Golang-internship.git/internal/config"
	"github.com/Masha003/Golang-internship.git/internal/data"
	"github.com/Masha003/Golang-internship.git/internal/data/repository"
	"github.com/Masha003/Golang-internship.git/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect postgres database")
	}

	rdb, err := data.NewRDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect redis database")
	}

	// User
	userRepository := repository.NewUserRepository(db, rdb)
	userService := service.NewUserService(userRepository, cfg)
	userController := controllers.NewUserController(userService)

	srv := api.NewServer(cfg, userController)

	//// Boilerpalte Gracefull Shutdown ////
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server")
		}
		log.Print("All server connections are closed")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown")
	}

	if err := data.CloseDB(db); err != nil {
		log.Fatal("Failed to close db connection")
	}

	if err := rdb.Close; err != nil {
		log.Fatal("Failed to close redis connection")
	}

	log.Print("Server exited properly")
}
