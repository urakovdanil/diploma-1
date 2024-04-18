package config

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log/slog"
)

const (
	defaultRunAddress           = "localhost:8080"
	defaultDatabaseURI          = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	defaultAccrualSystemAddress = "localhost:8081"
	defaultMigrationsFolder     = "./internal/migrations"
)

var (
	defaultLogLevel = slog.LevelInfo.String()
)

type StartUp struct {
	MigrationsFolder     string `env:"MIGRATIONS_FOLDER" json:"migrations_folder"`
	RunAddress           string `env:"RUN_ADDRESS" json:"run_address"`
	DatabaseURI          string `env:"DATABASE_URI" json:"database_uri"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" json:"accrual_system_address"`
	LogLevel             string `env:"LOG_LEVEL" json:"log_level"`
}

func (c *StartUp) String() string {
	js, err := json.Marshal(c)
	if err != nil {
		return "error on json.Marshal"
	}
	return string(js)
}

// New ...
func New(_ context.Context) (*StartUp, error) {
	res := &StartUp{
		RunAddress:           defaultRunAddress,
		DatabaseURI:          defaultDatabaseURI,
		AccrualSystemAddress: defaultAccrualSystemAddress,
		MigrationsFolder:     defaultMigrationsFolder,
	}

	if err := env.Parse(res); err != nil {
		return nil, err
	}

	ra := flag.String("a", res.RunAddress, fmt.Sprintf("address to run server on, defaults to %s", defaultRunAddress))
	du := flag.String("d", res.DatabaseURI, fmt.Sprintf("address to connect to PostgreSQL, defaults to %s", defaultDatabaseURI))
	asa := flag.String("r", res.AccrualSystemAddress, fmt.Sprintf("address to connect to accrual system, defaults to %s", defaultAccrualSystemAddress))
	ll := flag.String("l", defaultLogLevel, fmt.Sprintf("application log level, defaults to %s", defaultLogLevel))
	mf := flag.String("m", res.MigrationsFolder, fmt.Sprintf("path to migrations folder, defaults to %s", defaultMigrationsFolder))

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

	return res, nil
}
