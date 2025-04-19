package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"stan.com/stantest/config"
	"stan.com/stantest/routes"
)

func main() {
	// create a new echo instance
	e := echo.New()

	// set logging level if not existing with default debug level
	logLevelStr := os.Getenv("STAN_EPISODE_SERVER_LOG_LEVEL")
	logLevel := config.LOG_LEVEL_DEBUG
	if logLevelStr == "" {
		logLevel = log.Lvl(config.LOG_LEVEL_DEBUG)
	}
	e.Logger.SetLevel(logLevel)
	e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line}")

	// create or open the log file
	logFile, err := os.OpenFile("stantest.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		e.Logger.Fatal("failed to open log file:", err)
	}
	defer logFile.Close() // ensure log file closed correctly

	// bind the logger with our log file
	e.Logger.SetOutput(logFile)

	// add some default middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logFile,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// bind routes
	routes.SetupRoutes(e)

	// get API service port from system environment
	// if not setting , use default one
	port := os.Getenv("STAN_EPISODE_SERVER_PORT")
	if port == "" {
		port = config.DEFAULT_PORT
	}

	// start my server
	go func() {
		e.Logger.Infof("starting server on port %s", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	e.Logger.Info("received shutdown signal")

	// give some time to exit or shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown server safely now
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("server forced to shutdown:", err)
	}

	e.Logger.Info("server gracefully stopped")
}
