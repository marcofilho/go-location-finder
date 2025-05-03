package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/spf13/viper"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Request struct {
	Cep string `json:"cep"`
}

type Conf struct {
	API_KEY                     string `mapstructure:"API_KEY"`
	TITLE                       string `mapstructure:"TITLE"`
	CONTENT                     string `mapstructure:"CONTENT"`
	RESPONSE_TIME               string `mapstructure:"RESPONSE_TIME"`
	REQUEST_NAME_OTEL           string `mapstructure:"REQUEST_NAME_OTEL"`
	OTEL_SERVICE_NAME           string `mapstructure:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_ENDPOINT string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	HTTP_PORT                   string `mapstructure:"HTTP_PORT"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Weather struct {
	Location   Location `json:"location"`
	Current    Current  `json:"current"`
	WindMph    float64  `json:"wind_mph"`
	WindKph    float64  `json:"wind_kph"`
	WindDegree int      `json:"wind_degree"`
	WindDir    string   `json:"wind_dir"`
	PressureMb int      `json:"pressure_mb"`
	PressureIn float64  `json:"pressure_in"`
	PrecipMm   float64  `json:"precip_mm"`
	PrecipIn   float64  `json:"precip_in"`
	Humidity   int      `json:"humidity"`
	Cloud      int      `json:"cloud"`
	FeelslikeC float64  `json:"feelslike_c"`
	FeelslikeF float64  `json:"feelslike_f"`
	VisKm      float64  `json:"vis_km"`
	VisMiles   float64  `json:"vis_miles"`
	Uv         float64  `json:"uv"`
	GustMph    float64  `json:"gust_mph"`
	GustKph    float64  `json:"gust_kph"`
}

type Response struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

type Cep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	http.HandleFunc("/cep", cepHandler)
	port := ":8080"
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(port, nil)
}

func getCep(cep string) (Cep, error) {
	req, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return Cep{}, err
	}

	var c Cep
	err = json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		return Cep{}, err
	}
	defer req.Body.Close()

	return c, nil
}

func getWeather(city, apiKey string) (Weather, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)

	req, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	defer req.Body.Close()

	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var w Weather
	err = json.NewDecoder(req.Body).Decode(&w)
	if err != nil {
		return Weather{}, err
	}
	return w, nil
}

func convertToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func convertToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return output
}

func loadConfig() (*Conf, error) {
	var cfg Conf

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("API_KEY")
	viper.BindEnv("TITLE")
	viper.BindEnv("CONTENT")
	viper.BindEnv("RESPONSE_TIME")
	viper.BindEnv("REQUEST_NAME_OTEL")
	viper.BindEnv("OTEL_SERVICE_NAME")
	viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	viper.BindEnv("HTTP_PORT")

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Config loaded: %+v\n", cfg)

	return &cfg, err
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	var cep Request
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config, err := loadConfig()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	apiKey := config.API_KEY

	c, err := getCep(cep.Cep)
	if err != nil || c == (Cep{}) {
		http.Error(w, "Cannot find zipcode", http.StatusNotFound)
		return
	}

	type viacepError struct {
		Erro bool `json:"erro"`
	}
	var errCheck viacepError

	cepJson, _ := json.Marshal(c)
	json.Unmarshal(cepJson, &errCheck)
	if errCheck.Erro {
		http.Error(w, "Invalid zipcode", http.StatusBadRequest)
		return
	}

	location := removeAccents(c.Localidade)

	fmt.Println("Location: ", location)
	fmt.Println("API Key: ", apiKey)

	weather, err := getWeather(location, apiKey)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := Response{
		City:   weather.Location.Name,
		Temp_C: weather.Current.TempC,
		Temp_F: convertToFahrenheit(weather.Current.TempC),
		Temp_K: convertToKelvin(weather.Current.TempC),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
