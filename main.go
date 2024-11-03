package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/sysnote8main/readmyvc/internal/bot"
)

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	// Load dotenv
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	bot.Start()
}
