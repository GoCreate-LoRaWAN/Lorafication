// Package config exposes a Config type that pulls values from configuration files for
// lorafication daemon.
package config

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/george-e-shaw-iv/lorafication/internal/platform/duration"
)

// Constant block for Config struct field defaults.
const (
	// DefaultLogLevel is the default value of the LogLevel struct field on the Config
	// type.
	DefaultLogLevel = 0

	// DefaultPort is the default value of the Port struct field on the Config type.
	DefaultPort = 9000

	// DefaultReadTimeout is the default value of the ReadTimeout struct field on the
	// Config type.
	DefaultReadTimeout = 10 * time.Second

	// DefaultWriteTimeout is the default value of the WriteTimeout struct field on
	// the Config type.
	DefaultWriteTimeout = 20 * time.Second

	// DefaultShutdownTimeout is the default value of the ShutdownTimeout struct field
	// on the Config type.
	DefaultShutdownTimeout = 20 * time.Second
)

// Config is a struct that contains the struct fields necessary for running the
// lorafication daemon.
type Config struct {
	LogLevel        int               `json:"logLevel" yaml:"logLevel" env:"LOG_LEVEL"`
	Port            int               `json:"port" yaml:"port" env:"PORT"`
	ReadTimeout     duration.Duration `json:"readTimeout" yaml:"readTimeout" env:"READ_TIMEOUT"`
	WriteTimeout    duration.Duration `json:"writeTimeout" yaml:"writeTimeout" env:"WRITE_TIMEOUT"`
	ShutdownTimeout duration.Duration `json:"shutdownTimeout" yaml:"shutdownTimeout" env:"SHUTDOWN_TIMEOUT"`
}

// TODO(George): Create helper functions to set config from JSON/YAML/ENV using struct tags.

// Defaults is a method on the Config pointer receiver that sets defaults on the receiver
// where values are not already set.
func (c *Config) Defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = DefaultLogLevel
	}

	if c.Port == 0 {
		c.Port = DefaultPort
	}

	if c.ReadTimeout.IsEmpty() {
		c.ReadTimeout.Duration = DefaultReadTimeout
	}

	if c.WriteTimeout.IsEmpty() {
		c.WriteTimeout.Duration = DefaultWriteTimeout
	}

	if c.ShutdownTimeout.IsEmpty() {
		c.ShutdownTimeout.Duration = DefaultShutdownTimeout
	}
}

// Validate is a method on the Config pointer receiver that validates the values set on
// the receiver.
func (c *Config) Validate() error {
	if actual, min, max := zapcore.Level(c.LogLevel), zapcore.DebugLevel, zapcore.FatalLevel; actual < min || actual > max {
		return fmt.Errorf("log level must be [%d, %d]", min, max)
	}

	if c.Port <= 0 {
		return errors.New("port must be > 0")
	}

	if c.ReadTimeout.IsEmpty() {
		return errors.New("read timeout must be > 0ms")
	}

	if c.WriteTimeout.IsEmpty() {
		return errors.New("write timeout must be > 0ms")
	}

	if c.ShutdownTimeout.IsEmpty() {
		return errors.New("shutdown timeout must be > 0ms")
	}

	return nil
}
