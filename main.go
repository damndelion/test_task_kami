package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damndelion/test_task_kami/configs"
	_ "github.com/damndelion/test_task_kami/docs"
	"github.com/damndelion/test_task_kami/internal/handler"
	"github.com/damndelion/test_task_kami/internal/infrastructure/database"
	"github.com/damndelion/test_task_kami/internal/infrastructure/logger"
	"github.com/damndelion/test_task_kami/internal/infrastructure/server"
	"github.com/damndelion/test_task_kami/internal/repository"
	"github.com/damndelion/test_task_kami/internal/service"
)

//	@title			test_task_kami
//	@version		1.0
//	@description	API Server for test_task_kami

// @host		localhost:8080
// @BasePath	/
func main() {
	// Read configs
	config, err := configs.InitConfigs()
	if err != nil {
		log.Fatalf("error with loading config: %s", err.Error())
	}

	// Logger
	logs, err := logger.NewLogger(config.Logger)
	if err != nil {
		log.Fatalf("error with loading config: %s", err.Error())
	}

	// Connect to DB
	db, err := database.NewPostgresDB(config.Postgres)
	if err != nil {
		logs.Named("database").Fatalf("error with loading config: %s", err.Error())
		return
	}
	logs.Named("database").Info("connection with database is established")

	repo := repository.NewRepository(db)

	useCase := service.NewService(logs, repo)

	handlers := handler.NewHandler(logs, useCase)
	srv := new(server.Server)
	go func() {
		if err = srv.Run(config.Http.Port, handlers.InitRoutes()); err != nil {
			logs.Named("server").Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logs.Named("app").Infoln("App Started")

	// Implemented Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logs.Named("app").Infoln("App Shutting Down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logs.Named("server").Errorf("Error occurred during server shutdown: %s", err.Error())
	} else {
		logs.Named("server").Infoln("Server shutdown completed")
	}

	db.Close()
	logs.Named("db").Infoln("Database connection closed")

	select {
	case <-ctx.Done():
		logs.Named("app").Infoln("Shutdown timeout reached, exiting")
	}

	logs.Named("app").Infoln("Shutdown complete")
}
