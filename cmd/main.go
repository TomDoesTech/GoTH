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
	"github.com/tomdoestech/goth/internal/auth"
	"github.com/tomdoestech/goth/internal/pkg/config"
	"github.com/tomdoestech/goth/internal/pkg/metrics"
	users "github.com/tomdoestech/goth/internal/user"
	"github.com/tomdoestech/goth/internal/web"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tokenAuth *jwtauth.JWTAuth

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

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

	r.Use(metrics.NewPatternMiddleware(conf.ServiceName))

	tokenAuth = jwtauth.New("RS256", conf.JWTPrivateKey, conf.JWTPublicKey)

	r.Use(jwtauth.Verify(tokenAuth, TokenFromCookie))

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

	webHandler := web.NewWebHandler(
		web.WebHandlerParams{
			Logger: logger,
		},
	)

	auth.NewAuthHTTP(auth.AuthHTTPParams{
		AuthHandler: authHandler,
		Mux:         r,
	})

	web.NewWebHTTP(web.WebHTTPParams{
		WebHandler: webHandler,
		Mux:        r,
	})

	go metrics.StartMetricsServer(logger)

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
