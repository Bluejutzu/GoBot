package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bluejutzu/GoBot/handlers"
	"github.com/bwmarrin/discordgo"
)

var BotToken string

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.AddHandler(handlers.MessageCreate)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error as occured")
	}
}
