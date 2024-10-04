package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	CHECKER_SERVICE string
	DB_HOST      string
	DB_PORT      string
	DB_USER      string
	DB_PASSWORD  string
	DB_NAME      string
}

func LoadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("error loading .env file or not found", err)
	}

	config := Config{}

	config.CHECKER_SERVICE = cast.ToString(coalesce("CHECKER_SERVICE", ":50054"))
	config.DB_HOST = cast.ToString(coalesce("PDB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("PDB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("PDB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("PDB_PASSWORD", "1111"))
	config.DB_NAME = cast.ToString(coalesce("PDB_NAME", "testuzb1_checker"))
	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
