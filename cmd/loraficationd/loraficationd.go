package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/george-e-shaw-iv/lorafication/cmd/loraficationd/config"
	"github.com/george-e-shaw-iv/lorafication/cmd/loraficationd/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	f, err := os.Open("env.json")
	if err != nil {
		exitCode = 1
		log.Fatalf("open environment file: %v", err)
	}

	// Necessary because we don't want to defer closing this but need to handle
	// closing this in more than one place. The reason we don't want to defer this
	// is because if we do the resource will stick around for the entirety of the
	// program. Eventually this code will be abstracted into the config package.
	closer := func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("close environment file: %v", err)
		}
	}

	var cfg config.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		// Close environment file to free resource.
		closer(f)

		exitCode = 1
		log.Fatalf("parse env variables from file: %v", err)
	}

	// Close environment file to free resource.
	closer(f)

	zCfg := zap.NewProductionConfig()
	zCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.LogLevel))
	zCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zCfg.Build()
	if err != nil {
		exitCode = 1
		log.Fatalf("build logger: %v", err)
	}
	logger = logger.With(zap.String("service", "loraficationd"))

	cfg.Defaults()
	if err := cfg.Validate(); err != nil {
		logger.Error("validate config", zap.Error(err))
		exitCode = 1
		return
	}

	api := http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.Port),
		Handler:      server.NewServer(&cfg, logger),
		ReadTimeout:  cfg.ReadTimeout.Duration,
		WriteTimeout: cfg.WriteTimeout.Duration,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	apiErr := make(chan error, 1)
	go func() {
		logger.Info("server started", zap.Int("port", cfg.Port))
		apiErr <- api.ListenAndServe()
	}()

	select {
	case <-shutdown:
		logger.Info("shutdown signal received, attempting to gracefully terminate server")
		signal.Reset(os.Interrupt, syscall.SIGTERM)
	case err := <-apiErr:
		logger.Error("fatal server error", zap.Error(err))
		exitCode = 1
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout.Duration)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown error", zap.Error(err))
		logger.Info("attempting to forcefully shutdown server")

		if err := api.Close(); err != nil {
			logger.Error("forceful shutdown error", zap.Error(err))
			exitCode = 1
			return
		}
	}
}
