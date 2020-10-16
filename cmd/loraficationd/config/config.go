// Package config exposes a Config type that pulls values from configuration files for
// lorafication daemon.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/22arw/lorafication/internal/platform/duration"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

// Constant block for Config struct field defaults.
const (
	// DefaultLogLevel is the default value of the LogLevel struct field on the Config
	// type.
	DefaultLogLevel = 0

	// DefaultPort is the default value of the Port struct field on the Config type.
	DefaultPort = 9000

	// DefaultDBUser is the default value of the DBUser struct field on the Config type.
	DefaultDBUser = "root"

	// DefaultDBPass is the default value of the DBPass struct field on the Config type.
	DefaultDBPass = "root"

	// DefaultDBName is the default value of the DBName struct field on the Config type.
	DefaultDBName = "lorafication"

	// DefaultDBHost is the default value of the DBHost struct field on the Config type.
	DefaultDBHost = "db"

	// DefaultDBPort is the default value of the DBPort struct field on the Config type.
	DefaultDBPort = 5432

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
	LogLevel int `json:"logLevel" yaml:"logLevel" envconfig:"LOG_LEVEL"`
	Port     int `json:"port" yaml:"port" envconfig:"PORT"`

	DBUser string `json:"dbUser" yaml:"dbUser" envconfig:"DB_USER"`
	DBPass string `json:"dbPass" yaml:"dbPass" envconfig:"DB_PASS"`
	DBName string `json:"dbName" yaml:"dbName" envconfig:"DB_NAME"`
	DBHost string `json:"dbHost" yaml:"dbHost" envconfig:"DB_HOST"`
	DBPort int    `json:"dbPort" yaml:"dbPort" envconfig:"DB_PORT"`

	ReadTimeout     duration.Duration `json:"readTimeout" yaml:"readTimeout" envconfig:"READ_TIMEOUT"`
	WriteTimeout    duration.Duration `json:"writeTimeout" yaml:"writeTimeout" envconfig:"WRITE_TIMEOUT"`
	ShutdownTimeout duration.Duration `json:"shutdownTimeout" yaml:"shutdownTimeout" envconfig:"SHUTDOWN_TIMEOUT"`
}

// FromEnvironment gathers the configuration variables from the environment.
func FromEnvironment() (Config, error) {
	var c Config

	if err := envconfig.Process("LORAFICATION", &c); err != nil {
		return c, fmt.Errorf("process environment variables: %w", err)
	}
	c.Defaults()

	if err := c.Validate(); err != nil {
		return c, fmt.Errorf("validate configuration: %w", err)
	}

	return c, nil
}

// FromFile gathers the configuration variables from either a JSON or YAML file.
func FromFile(fp string) (Config, error) {
	var c Config

	f, err := os.Open(fp)
	if err != nil {
		return c, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	switch ext := filepath.Ext(fp); ext {
	case ".json":
		if err := json.NewDecoder(f).Decode(&c); err != nil {
			return c, fmt.Errorf("decode json file: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.NewDecoder(f).Decode(&c); err != nil {
			return c, fmt.Errorf("decode json file: %w", err)
		}
	}
	c.Defaults()

	if err := c.Validate(); err != nil {
		return c, fmt.Errorf("validate configuration: %w", err)
	}

	return c, nil
}

// Defaults is a method on the Config pointer receiver that sets defaults on the receiver
// where values are not already set.
func (c *Config) Defaults() {
	if c.LogLevel == 0 {
		c.LogLevel = DefaultLogLevel
	}

	if c.Port == 0 {
		c.Port = DefaultPort
	}

	if c.DBUser == "" {
		c.DBUser = DefaultDBUser
	}

	if c.DBPass == "" {
		c.DBPass = DefaultDBPass
	}

	if c.DBName == "" {
		c.DBName = DefaultDBName
	}

	if c.DBHost == "" {
		c.DBHost = DefaultDBHost
	}

	if c.DBPort == 0 {
		c.DBPort = DefaultDBPort
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

	if c.DBUser == "" {
		return errors.New("db user must be defined")
	}

	if c.DBPass == "" {
		return errors.New("db pass must be defined")
	}

	if c.DBName == "" {
		return errors.New("db name must be defined")
	}

	if c.DBHost == "" {
		return errors.New("db host must be defined")
	}

	if c.DBPort <= 0 {
		return errors.New("db port must be > 0 ")
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
