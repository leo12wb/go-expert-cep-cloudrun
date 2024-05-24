package get_temperature

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	zipcode2 "github.com/leo12wb/go-expert-cep-cloudrun/internal/entity/zipcode"
	"github.com/stretchr/testify/suite"
)

type GetTemperatureTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *http.Client
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(GetTemperatureTestSuite))
}

func (suite *GetTemperatureTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.client = &http.Client{}
}

// Test case for getCityFromZipCode function using mock server
func (suite *GetTemperatureTestSuite) TestGetCityFromZipCode() {
	// Mock API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AddressResponse{Localidade: "Dublin"}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	city, err := getCityFromZipCode(mockServer.URL+"/%s/json/", "12345")

	// Check if there's no error
	suite.NoError(err)
	suite.Equal("Dublin", city)
}

// Test case for getTemperatureFromCity function using mock server
func (suite *GetTemperatureTestSuite) TestGetTemperatureFromCity() {
	// Mock API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"current": map[string]interface{}{"temp_c": 10.5}}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockServer.Close()

	// Call the function with the mock server URL
	temperature, err := getTemperatureFromCity(mockServer.URL+"/current.json?key=%s&q=%s", "your-api-key", "Dublin")

	// Check if there's no error
	suite.NoError(err)
	suite.Equal(10.5, temperature)
}
func (suite *GetTemperatureTestSuite) TestUseCase_Execute() {
	// Mock API responses
	mockApiCepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := AddressResponse{Localidade: "Dublin"}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockApiCepServer.Close()

	mockWeatherApiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{"current": map[string]interface{}{"temp_c": 10.5}}
		jsonBytes, _ := json.Marshal(response)
		_, _ = w.Write(jsonBytes)
	}))
	defer mockWeatherApiServer.Close()

	// Create a test inputDTO with the URLs of the mock servers
	zipcode, err := zipcode2.NewZipcode("12345678")

	input := InputDTO{
		WeatherApiKey: "your-api-key",
		WeatherApiUrl: mockWeatherApiServer.URL + "/current.json?key=%s&q=%s",
		ApiCepUrl:     mockApiCepServer.URL + "/%s/json/",
		Zipcode:       zipcode, // replace with a valid zip code for testing
	}
	// Create the use case instance with test inputDTO, context and http client
	uc := NewUseCase(&suite.ctx, suite.client)

	// Execute the use case
	output, err := uc.Execute(input)

	// Check if there's no error
	suite.NoError(err)

	// Check if city is not empty
	suite.NotEmpty(output.City)

	// Check if temperatures are not zero
	suite.NotZero(output.CelsiusTemp)
	suite.NotZero(output.FahrenheitTemp)
	suite.NotZero(output.KelvinTemp)
}
