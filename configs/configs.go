package configs

import (
	"messanger/internal/logs"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logs.FatalLog("", "No .env file found", nil)
	}
}
