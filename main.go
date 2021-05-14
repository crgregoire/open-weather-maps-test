package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//
//OpenWeather : Struct for the whole response body
//
type OpenWeather struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

//
//Coord : struct for coordinates
//
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

//
// Weather : struct for weather title, description, and icon
//
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

//
// Main : struct for temperatures
//
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

//
// Wind : struct for wind info
//
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

//
// Clouds : struct for cloud info
//
type Clouds struct {
	All int `json:"all"`
}

//
// Sys : struct for system info like country
//
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

//
// main function handles all input and response from the API
//
func main() {
	var latitude float64
	var longitude float64
	fmt.Print("Enter your latitude: ")

	_, err := fmt.Scanf("%f", &latitude)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter your longitude: ")

	_, err = fmt.Scanf("%f", &longitude)

	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=" + fmt.Sprintf("%f", latitude) + "&lon=" + fmt.Sprintf("%f", longitude) + "&units=imperial&appid=67bef9f0f5d501c8d4404d04d1d2ba2f")

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\nGetting your location's weather data!\n\n")

	var weatherJSON OpenWeather
	err = json.Unmarshal(responseData, &weatherJSON)
	if err != nil {
		log.Fatal(err)
	}

	printTemp(weatherJSON)
	temperatureDesc(weatherJSON)
	returnConditionAndAlerts(weatherJSON)

}

//
// returnConditionAndAlerts function handles wind direction and speed, checks for an alert (any weather Description
//	under 800 as I don't have access to the paid Alerts API), and returns conditions.
//
func returnConditionAndAlerts(weatherData OpenWeather) {

	var directions = []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	i := int((float64(weatherData.Wind.Deg) + 11.25) / 22.5)

	fmt.Println("The weather condition outside is: " + weatherData.Weather[0].Description)
	fmt.Printf("The wind direction and speed is %.2fmph "+directions[i%16]+"\n", weatherData.Wind.Speed)

	alertCheck := float64(weatherData.Weather[0].ID) / 200

	if alertCheck < 4 {
		fmt.Println("ALERT! THERE IS A " + strings.ToUpper(weatherData.Weather[0].Description) + " OUTSIDE. " +
			"\nAlso alerts cost $125 a month so this is what you get instead sorry." +
			"\nBut if it was free it's just another JSON object")
	}

	if float64(weatherData.Visibility) < 805 {
		fmt.Println("ALERT: LOW VISIBILITY")
	}
}

//
// temperatureDesc function describes if it's cold, okay, warm, or hot.
//
func temperatureDesc(weatherData OpenWeather) {

	if weatherData.Main.Temp < 0 {
		fmt.Print("Oof yeah stay inside or you know, you might die.\n\n")
	} else if weatherData.Main.Temp >= 0 && weatherData.Main.Temp < 50 {
		fmt.Print("Brrr... It's a bit cold outside.\n\n")
	} else if weatherData.Main.Temp >= 50 && weatherData.Main.Temp < 85 {
		fmt.Print("The temperature outside is just right.\n\n")
	} else if weatherData.Main.Temp >= 85 && weatherData.Main.Temp < 100 {
		fmt.Print("Yeah it's a bit warm out.\n\n")
	} else if weatherData.Main.Temp > 100 {
		fmt.Print("Alright yeah it's hot, where is the AC?\n\n")
	} else {
		fmt.Print("It's something outside. Idk?\n\n")
	}
}

//
// printTemp function literally just prints the temperature
//
func printTemp(weatherData OpenWeather) {
	fmt.Println("You're location is " + weatherData.Name)
	fmt.Printf("The temperature outside is %.2fF \n", weatherData.Main.Temp)
}
