package config

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

const (
	defaultRunAddress           = "localhost:8080"
	defaultDatabaseURI          = "postgres:5432"
	defaultAccrualSystemAddress = "localhost:8081"
)

type StartUp struct {
	RunAddress           string `env:"RUN_ADDRESS" json:"run_address"`
	DatabaseURI          string `env:"DATABASE_URI" json:"database_uri"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" json:"accrual_system_address"`
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
	}

	if err := env.Parse(res); err != nil {
		return nil, err
	}

	ra := flag.String("a", res.RunAddress, fmt.Sprintf("address to run server on, defaults to %s", defaultRunAddress))
	du := flag.String("d", res.DatabaseURI, fmt.Sprintf("address to connect to PostgreSQL, defaults to %s", defaultDatabaseURI))
	asa := flag.String("r", res.AccrualSystemAddress, fmt.Sprintf("address to connect to accrual system, defaults to %s", defaultAccrualSystemAddress))

	if res.RunAddress == "" {
		res.RunAddress = *ra
	}
	if res.DatabaseURI == "" {
		res.DatabaseURI = *du
	}
	if res.AccrualSystemAddress == "" {
		res.AccrualSystemAddress = *asa
	}

	return res, nil
}
