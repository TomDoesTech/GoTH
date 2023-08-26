package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomdoestech/goth/internal/auth"
	"github.com/tomdoestech/goth/internal/middleware"
	"github.com/tomdoestech/goth/internal/pkg/config"
	users "github.com/tomdoestech/goth/internal/user"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func startMetricsServer(logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
	logger.Info("Metrics server started on port 9100")
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {

	validate = validator.New()

	conf := config.Must()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sugar := logger.Sugar()

	r := chi.NewRouter()

	r.Use(middleware.RenderMiddleware)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	tokenAuth = jwtauth.New("RS256", conf.JWTPrivateKey, conf.JWTPublicKey)

	usersService := users.NewUserService(users.UserServiceParams{
		Logger:   logger,
		Validate: validate,
		DB:       db,
	})
	authService := auth.NewAuthService(auth.AuthServiceParams{
		Logger:    logger,
		SecretKey: []byte("secret"),
		TokenAuth: tokenAuth,
	})

	authHandler := auth.NewAuthHandler(
		auth.AuthHandlerParams{
			AuthService: authService,
			UserService: usersService,
			Validate:    validate,
			Logger:      logger,
		},
	)

	auth.NewAuthHTTP(auth.AuthHTTPParams{
		AuthHandler: authHandler,
		Mux:         r,
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":  "My Website",
			"Header": "Welcome to my website!",
			"Footer": "© 2023 My Website",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "home.html", data)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":  "About",
			"Header": "Welcome to my website!",
			"Footer": "© 2023 My Website",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "about.html", data)
	})

	go startMetricsServer(logger)

	srv := &http.Server{
		Addr:    conf.Port,
		Handler: r,
	}

	// Listen for OS signals to initiate graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		sugar.Info(context.Background(), "Starting server on port %s", conf.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {

			sugar.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop // Block until a signal is received

	log.Println("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalln("Error shutting down server", zap.Error(err))
	}

	log.Println("Server gracefully stopped")
}
