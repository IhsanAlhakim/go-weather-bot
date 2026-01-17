package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/IhsanAlhakim/go-weather-bot/pkg/command"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	botToken := os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New("Bot " + botToken)

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("BOT READY")
	})

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := command.Handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	discord.Open()
	defer discord.Close()

	log.Println("Adding Commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(command.List))
	for i, v := range command.List {
		command, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = command
	}

	log.Println("BOT RUNNING...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

}
