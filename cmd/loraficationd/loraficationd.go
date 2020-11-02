package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/22arw/lorafication/cmd/loraficationd/config"
	"github.com/22arw/lorafication/cmd/loraficationd/server"
	"github.com/22arw/lorafication/internal/mail"
	"github.com/22arw/lorafication/internal/platform/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// flagEnvFile is a variable that is collected from the -envFile="" flag upon running the
// lorafication binary.
var flagEnvFile string

// flagVerbose is a variable that is collected from the -verbose flag upon running the
// lorafication binary.
var flagVerbose bool

func init() {
	// Register the collection of the -envFile="" flag if it was passed to the binary.
	flag.StringVar(&flagEnvFile, "envFile", "", "the path of a JSON or YAML file that contains configuration variables, if not supplied the configuration will be collected from the environment.")

	// Register the collection of the -verbose flagif it was passed to the binary.
	flag.BoolVar(&flagVerbose, "verbose", false, "display verbose information")

	// Collect all registered CLI flags.
	flag.Parse()
}

// main is the entrypoint for the lorafication daemon.
func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	var cfg config.Config
	var err error

	// Process configuration variables from either a file or the environment
	// depending on whether or not the -envFile="" flag was specified.
	if flagEnvFile != "" {
		if cfg, err = config.FromFile(flagEnvFile); err != nil {
			exitCode = 1
			log.Fatalf("collect config from file: %v", err)
		}
	} else {
		if cfg, err = config.FromEnvironment(); err != nil {
			exitCode = 1
			log.Fatalf("collect config from environment: %v", err)
		}
	}

	// Configure the logger.
	zCfg := zap.NewProductionConfig()
	zCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.LogLevel))
	zCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Build the logger.
	logger, err := zCfg.Build()
	if err != nil {
		exitCode = 1
		log.Fatalf("build logger: %v", err)
	}

	// Add a service key with the value of loraficationd to all logs emitted
	// from this daemon to namespace the logs.
	logger = logger.With(zap.String("service", "loraficationd"))

	// Display verbose information if it was requested.
	if flagVerbose {
		// Log the parsed CLI flags.
		logger.Info("values of CLI flags",
			zap.String("envFile", flagEnvFile),
			zap.Bool("verbose", flagVerbose))

		// Log the parsed configuration values.
		logger.Info("values of configuration",
			zap.Int("port", cfg.Port),
			zap.Int("logLevel", cfg.LogLevel),
			zap.String("dbUser", cfg.DBUser),
			zap.String("dbName", cfg.DBName),
			zap.String("dbHost", cfg.DBHost),
			zap.Int("dbPort", cfg.DBPort),
			zap.String("smtpHost", cfg.SMTPHost),
			zap.Int("smtpPort", cfg.SMTPPort),
			zap.String("smtpUser", cfg.SMTPUser),
			zap.Duration("readTimeout", cfg.ReadTimeout.Duration),
			zap.Duration("writeTimeout", cfg.WriteTimeout.Duration),
			zap.Duration("shutdownTimeout", cfg.ShutdownTimeout.Duration))
	}

	// Construct database configuration struct to pass to the connection method.
	dbCfg := db.Config{
		User: cfg.DBUser,
		Pass: cfg.DBPass,
		Name: cfg.DBName,
		Host: cfg.DBHost,
		Port: cfg.DBPort,
	}

	// Connect to the database.
	dbc, err := db.NewConnection(logger, dbCfg)
	if err != nil {
		logger.Error("connect to database", zap.Error(err))
		exitCode = 1
		return
	}

	// Defer the closing of the database connection until after func main returns.
	defer func() {
		if err := dbc.Close(); err != nil {
			logger.Warn("error closing database connection", zap.Error(err))
		}
	}()

	// Configure the mailer used to send emails over SMTP.
	mailer := mail.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	// Configure the HTTP server that this daemon will expose.
	api := http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.Port),
		Handler:      server.NewServer(&cfg, logger, dbc, mailer),
		ReadTimeout:  cfg.ReadTimeout.Duration,
		WriteTimeout: cfg.WriteTimeout.Duration,
	}

	// Create a channel for interrupt and termination signals to be caught on
	// to possibly facilitate graceful shutdowns.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Create a channel to catch API-related non-recoverable errors.
	apiErr := make(chan error, 1)

	// Start the HTTP server.
	go func() {
		logger.Info("server started", zap.Int("port", cfg.Port))
		apiErr <- api.ListenAndServe()
	}()

	// Block until either a shutdown signal or an API-related non-recoverable error
	// is encountered.
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

	// Attempt to gracefully shutdown.
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
