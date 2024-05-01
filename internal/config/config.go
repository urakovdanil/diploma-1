package config

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log/slog"
	"sync"
)

const (
	defaultRunAddress           = "localhost:8080"
	defaultDatabaseURI          = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	defaultAccrualSystemAddress = "localhost:8081"
	defaultMigrationsFolder     = "./internal/migrations"
	defaultJWTTokenTTLMinutes   = 60
)

var (
	defaultLogLevel = slog.LevelInfo.String()
	Applied         *StartUp
)

type StartUp struct {
	MigrationsFolder     string `env:"MIGRATIONS_FOLDER" json:"migrations_folder"`
	RunAddress           string `env:"RUN_ADDRESS" json:"run_address"`
	DatabaseURI          string `env:"DATABASE_URI" json:"database_uri"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" json:"accrual_system_address"`
	LogLevel             string `env:"LOG_LEVEL" json:"log_level"`
	JWTSecretKey         string `env:"JWT_SECRET_KEY"`
	JWTTokenTTLMinutes   int    `env:"JWT_TOKEN_TTL_MINUTES" json:"jwt_token_ttl_minutes"`
	mu                   sync.RWMutex
}

func (c *StartUp) String() string {
	js, err := json.Marshal(c)
	if err != nil {
		return "error on json.Marshal"
	}
	return string(js)
}

func (c *StartUp) GetJWTSecretKey() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return []byte(c.JWTSecretKey)
}

func (c *StartUp) GetJWTTokenTTLMinutes() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.JWTTokenTTLMinutes
}

func (c *StartUp) GetLogLevel() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.LogLevel
}

func (c *StartUp) GetRunAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.RunAddress
}

func (c *StartUp) GetDatabaseURI() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.DatabaseURI
}

func (c *StartUp) GetAccrualSystemAddress() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AccrualSystemAddress
}

func (c *StartUp) GetMigrationsFolder() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.MigrationsFolder
}

// New ...
func New(_ context.Context) error {
	res := &StartUp{
		RunAddress:           defaultRunAddress,
		DatabaseURI:          defaultDatabaseURI,
		AccrualSystemAddress: defaultAccrualSystemAddress,
		MigrationsFolder:     defaultMigrationsFolder,
		JWTTokenTTLMinutes:   defaultJWTTokenTTLMinutes,
	}

	if err := env.Parse(res); err != nil {
		return err
	}

	ra := flag.String("a", res.RunAddress, fmt.Sprintf("address to run server on, defaults to %s", defaultRunAddress))
	du := flag.String("d", res.DatabaseURI, fmt.Sprintf("address to connect to PostgreSQL, defaults to %s", defaultDatabaseURI))
	asa := flag.String("r", res.AccrualSystemAddress, fmt.Sprintf("address to connect to accrual system, defaults to %s", defaultAccrualSystemAddress))
	ll := flag.String("l", defaultLogLevel, fmt.Sprintf("application log level, defaults to %s", defaultLogLevel))
	mf := flag.String("m", res.MigrationsFolder, fmt.Sprintf("path to migrations folder, defaults to %s", defaultMigrationsFolder))
	jwtTTL := flag.Int("j", res.JWTTokenTTLMinutes, fmt.Sprintf("jwt token lifetime in minutes, defaults to %d", defaultJWTTokenTTLMinutes))
	flag.Parse()

	if res.RunAddress == "" {
		res.RunAddress = *ra
	}
	if res.DatabaseURI == "" {
		res.DatabaseURI = *du
	}
	if res.AccrualSystemAddress == "" {
		res.AccrualSystemAddress = *asa
	}
	if res.LogLevel == "" {
		res.LogLevel = *ll
	}
	if res.MigrationsFolder == "" {
		res.MigrationsFolder = *mf
	}
	if res.JWTTokenTTLMinutes == 0 {
		res.JWTTokenTTLMinutes = *jwtTTL
	}

	Applied = res

	return nil
}
