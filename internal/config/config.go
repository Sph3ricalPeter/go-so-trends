package config

import (
	"os"
)

var (
	HOST    = os.Getenv("HOST")
	PORT    = os.Getenv("PORT")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_PASS = os.Getenv("DB_PSW")
)
