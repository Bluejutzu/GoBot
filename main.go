package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/bluejutzu/GoBot/bot"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botToken := os.Getenv("BOT_TOKEN")
	bot.BotToken = botToken
	bot.Run()
}
