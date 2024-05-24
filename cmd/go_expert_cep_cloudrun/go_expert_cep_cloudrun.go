package go_expert_cep_cloudrun

import (
	"context"
	"net/http"
	"os"

	"github.com/leo12wb/go-expert-cep-cloudrun/cmd/go_expert_cep_cloudrun/dependency_injection"
	"github.com/leo12wb/go-expert-cep-cloudrun/configs"
	"github.com/leo12wb/go-expert-cep-cloudrun/internal/infra/web/webserver"
	"github.com/rs/zerolog/log"
)

func handleErr(err error) {
	if err != nil {
		log.Error().Err(err).Msg("")
		panic(err)
	}
}

func Bootstap() {
	workdir, err := os.Getwd()
	handleErr(err)
	appConfig, err := configs.LoadConfig(workdir)
	if err != nil {
		panic(err)
	}

	restServer := webserver.NewWebServer(appConfig.WebserverPort)
	ctx := context.Background()
	temperatureHandler := dependency_injection.NewTemperatureHandler(&ctx)

	temperatureHandler.WeatherApiKey = appConfig.WeatherApiKey
	temperatureHandler.ApiCepUrl = appConfig.CepApiURL
	temperatureHandler.WeatherApiUrl = appConfig.WeatherApiURL

	restServer.AddHandler("/", http.MethodGet, temperatureHandler.Handle)
	restServer.Start()

}
