package main

import "github.com/ObiaNzk/LTK-test-manu/internal/platform"

type config struct {
	IsProduction bool
	DBConfig     platform.DBConfig
}

func newLocalConfig() config {
	dbConfig := platform.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "events_db",
		SSLMode:  "disable",
	}

	config := config{
		IsProduction: false,
		DBConfig:     dbConfig,
	}

	return config

}
