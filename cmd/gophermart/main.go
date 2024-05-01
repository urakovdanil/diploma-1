package main

import (
	"context"
	"diploma-1/internal/api/auth"
	middleware "diploma-1/internal/api/middleware"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// TODO: добавить middleware с logger.With() для добавления трейса запроса

func main() {
	ctx := context.Background()
	if err := config.New(ctx); err != nil {
		fmt.Printf("unable to collect config: %v", err)
	}
	fmt.Printf("applied args: %s\n", config.Applied)
	logger.New(config.Applied)
	if err := storage.New(ctx, config.Applied); err != nil {
		logger.Fatalf(ctx, "unable to init storage: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.CustomizeResponseWriter)
	router.Use(middleware.ResponseCompressor)
	router.Use(middleware.RequestDecompressor)
	router.Use(middleware.RequestLogger)
	router.Use(chiMiddleware.Recoverer)

	authRouter := auth.New()
	router.Post(auth.RegisterPath, authRouter.RegisterHandlerFunc)
	router.Post(auth.LoginPath, authRouter.LoginHandlerFunc)

	srv := &http.Server{
		Addr:    config.Applied.GetRunAddress(),
		Handler: router,
	}
	go func() {
		logger.Infof(ctx, "starting server on %s", config.Applied.GetRunAddress())
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info(ctx, "server closed")
			} else {
				logger.Errorf(ctx, "unable to start server: %v", err)
			}
		}
	}()

	sd := make(chan os.Signal, 1)
	signal.Notify(sd, os.Interrupt)
	<-sd
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf(ctx, "unable to shutdown server: %v", err)
	}
	// TODO: implement closer
}
