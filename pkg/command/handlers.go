package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type CityWeather struct {
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Humidity   int     `json:"humidity"`
	} `json:"main"`
	Name    string `json:"name"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Cod int `json:"cod"`
}

var Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, option := range options {
			optionMap[option.Name] = option
		}

		var city string

		if option, ok := optionMap["city"]; ok {
			city = option.StringValue()
		}

		if city == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please input a city name",
				},
			})
			return
		}

		var client = &http.Client{}
		var data CityWeather

		openWeatherKey := os.Getenv("OPEN_WEATHER_API_KEY")

		var url = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%v&units=metric&appid=%v", city, openWeatherKey)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error Server",
				},
			})
			return
		}

		response, err := client.Do(request)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error Server",
				},
			})
			return
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "City not found",
				},
			})
			return
		}

		err = json.NewDecoder(response.Body).Decode(&data)

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error Server",
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("City: %v, Weather: %v, Description: %v, Temperature: %v", data.Name, data.Weather[0].Main, data.Weather[0].Description, data.Main.Temp),
			},
		})
	},
}
