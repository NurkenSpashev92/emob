package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nurkenspashev92/emob/cmd/router"
	"github.com/nurkenspashev92/emob/configs"
	"github.com/nurkenspashev92/emob/pkg/store"
)

type App struct {
	fiberApp *fiber.App
}

func (a *App) Run() {
	cfg := configs.NewConfig()

	database, err := store.NewPostgresDb(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer database.Close()

	a.fiberApp = router.RegisterRoutes(database.Conn)
	done := make(chan bool, 1)
	go func() {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		log.Printf("Server running on %d...\n", port)

		err := a.fiberApp.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go a.Shutdown(done)
	<-done
	log.Println("Graceful shutdown complete.")
}

func (fiberServer *App) Shutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.fiberApp.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}
