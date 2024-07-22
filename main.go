package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sysnote8main/readmyvc/script"
)

func main() {
	// Process dotenv
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	// Discord bot
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Failed to start discord session", err)
		return
	}

	dg.AddHandler(script.OnMessageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Failed to open connection", err)
		return
	}

	defer dg.Close()

	fmt.Println("Successfully to login!")

	// Bot stop logic
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}
