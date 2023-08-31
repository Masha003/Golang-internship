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
		log.Fatal("Failed to load config: ", err)
	}

	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect postgres database: ", err)
	}

	rdb, err := data.NewRDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect redis database: ", err)
	}

	// User
	userRepository := repository.NewUserRepository(db, rdb)
	userService := service.NewUserService(userRepository, cfg)
	userController := controllers.NewUserController(userService)

	// Start HTTP Server
	httpSrv := api.NewHttpServer(cfg, userController)
	go func() {
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server: ", err)
		}
		log.Print("All server connections are closed")
	}()

	// Start GRPC Server
	grpcSrv, listener, err := api.NewGrpcServer(cfg, userService)
	if err != nil {
		log.Fatal("Failed to create grpc server: ", err)
	}
	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatal("Grpc server failed to start: ", err)
		}
	}()

	// Gracefull Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
	log.Print("Shutting down server...")

	// Shutdown HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	// Shutdown GRPC Server
	grpcSrv.GracefulStop()

	// Close Postgress Connection
	if err := data.CloseDB(db); err != nil {
		log.Fatal("Failed to close db connection: ", err)
	}

	// Close Redis Connection
	if err := rdb.Close(); err != nil {
		log.Fatal("Failed to close redis connection: ", err)
	}

	log.Print("Server exited properly")
}
