package command

import "github.com/bwmarrin/discordgo"

var List = []*discordgo.ApplicationCommand{
	{
		Name:        "weather",
		Description: "check weather in a city",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "city",
				Description: "Enter a city that you want to know its weather",
				Required:    true,
			},
		},
	},
}
