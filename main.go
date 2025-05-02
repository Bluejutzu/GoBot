package main

import (
	"log"
	"os"

	"github.com/bluejutzu/GoBot/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load configuration from environment variables
	bot.BotToken = os.Getenv("BOT_TOKEN")
	bot.Run()
}
