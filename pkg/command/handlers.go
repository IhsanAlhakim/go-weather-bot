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

var serverErrorMessage = "Something went wrong"

var Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		if err := deferInteractionResponse(s, i); err != nil {
			return
		}

		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, option := range options {
			optionMap[option.Name] = option
		}

		var city string

		if option, ok := optionMap["city"]; ok {
			city = option.StringValue()
		}

		var client = &http.Client{}
		var data CityWeather

		openWeatherKey := os.Getenv("OPEN_WEATHER_API_KEY")

		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%v&units=metric&appid=%v", city, openWeatherKey)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			singleTextResponse(s, i, &serverErrorMessage)
			return
		}

		response, err := client.Do(request)
		if err != nil {
			singleTextResponse(s, i, &serverErrorMessage)
			return
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			singleTextResponse(s, i, &serverErrorMessage)
			return
		}

		err = json.NewDecoder(response.Body).Decode(&data)

		if err != nil {
			singleTextResponse(s, i, &serverErrorMessage)
			return
		}

		responseText := fmt.Sprintf("Weather Info\nCity: %v\nWeather: %v\nDescription: %v\nTemperature: %v", data.Name, data.Weather[0].Main, data.Weather[0].Description, data.Main.Temp)
		singleTextResponse(s, i, &responseText)
	},
}
