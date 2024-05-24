package get_temperature

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/leo12wb/go-expert-cep-cloudrun/internal/entity/zipcode"
)

type AddressResponse struct {
	Localidade string `json:"localidade"`
}

type InputDTO struct {
	WeatherApiKey string
	WeatherApiUrl string
	ApiCepUrl     string
	Zipcode       zipcode.Zipcode
}
type OutputDTO struct {
	City           string  `json:"-"`
	CelsiusTemp    float64 `json:"temp_C"`
	FahrenheitTemp float64 `json:"temp_F"`
	KelvinTemp     float64 `json:"temp_K"`
}

type UseCase struct {
	ctx    *context.Context
	client *http.Client
}

func NewUseCase(
	ctx *context.Context,
	client *http.Client,
) UseCase {
	return UseCase{
		ctx:    ctx,
		client: client,
	}
}

func getCityFromZipCode(apiUrl string, zipcode string) (string, error) {
	// "https://viacep.com.br/ws/%s/json/"
	resp, err := http.Get(fmt.Sprintf(apiUrl, zipcode))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var addressResponse AddressResponse
	err = json.Unmarshal(body, &addressResponse)
	if err != nil {
		return "", err
	}

	return addressResponse.Localidade, nil
}

func getTemperatureFromCity(apiUrl string, apikey string, city string) (float64, error) {
	//"http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
	url := fmt.Sprintf(apiUrl, url.QueryEscape(apikey), url.QueryEscape(city))
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var weatherResponse map[string]interface{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return 0, err
	}

	temperatureData, ok := weatherResponse["current"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("Invalid temperature data")
	}

	return temperatureData["temp_c"].(float64), nil
}

func (uc *UseCase) Execute(inputDto InputDTO) (OutputDTO, error) {
	var err error
	city, err := getCityFromZipCode(inputDto.ApiCepUrl, inputDto.Zipcode.Zipcode)
	result := OutputDTO{}
	if err != nil {
		return result, err
	}
	if city != "" {
		celsiusTemp, err := getTemperatureFromCity(inputDto.WeatherApiUrl, inputDto.WeatherApiKey, city)
		if err != nil {
			return result, err
		}
		result.CelsiusTemp = celsiusTemp
	}

	result.City = city
	result.FahrenheitTemp = (result.CelsiusTemp * 1.8) + 32
	result.KelvinTemp = result.CelsiusTemp + 273
	return result, nil
}
