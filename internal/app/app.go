package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taans/internal/telegram"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Config struct {
	Port string
}

func LoadConfig() (Config, error) {
	var config Config

	port, found := os.LookupEnv("PORT")
	if !found {
		fmt.Println(found)
		return config, fmt.Errorf("PORT not found")
	}
	config.Port = port

	return config, nil
}

type Application struct {
	Router   *chi.Mux
	Telegram telegram.Config
	Config   Config
	Message  chan telegram.Message
}

// Loads configuration from environment variables and creates a new Application
func NewApplication() (*Application, error) {
	t, err := telegram.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	router := chi.NewRouter()

	// Add middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"*"}, // Allow all methods
		AllowedHeaders:   []string{"*"}, // Allow all headers
		AllowCredentials: true,
	}))

	app := &Application{
		Telegram: t,
		Config:   c,
		Router:   router,
		Message:  make(chan telegram.Message),
	}

	return app, nil
}

func (app *Application) Start() error {
	go telegram.StartBot(app.Telegram, app.Message)

	server := &http.Server{
		Addr:         ":" + app.Config.Port,
		Handler:      app.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		log.Println("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}
