package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	HOST    string
	PORT    string
	DB_HOST string
	DB_PORT string
	DB_PASS string
)

func init() {
	godotenv.Load()

	HOST = os.Getenv("HOST")
	PORT = os.Getenv("PORT")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_PASS = os.Getenv("DB_PSW")
}
