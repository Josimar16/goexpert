package usecase_test

import (
	"testing"

	"github.com/josimar16/goexpert/desafio-cloud-run/internal/entity"
	"github.com/josimar16/goexpert/desafio-cloud-run/internal/infra/web"
	"github.com/josimar16/goexpert/desafio-cloud-run/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateCEP(t *testing.T) {
	get_weather_dto := usecase.ValidateCEPInputDTO{
		CEP: "01001001",
	}

	mockCEPRepo := new(MockCEPRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockCEPRepo.On("IsValid", get_weather_dto.CEP).Return(true)

	mockCEPRepo.On("Get", get_weather_dto.CEP).Return([]byte(`{
	        "cep": "01001-001",
	        "logradouro": "Praça da Sé",
	        "complemento": "lado par",
	        "bairro": "Sé",
	        "localidade": "São Paulo",
	        "uf": "SP",
	        "ibge": "3550308",
	        "gia": "1004",
	        "ddd": "11",
	        "siafi": "7107"
	    }`), nil)

	mockCEPRepo.On("Convert", mock.Anything).Return(&entity.CEP{
		CEP:         "01001-001",
		Logradouro:  "Praça da Sé",
		Complemento: "lado par",
		Bairro:      "Sé",
		Localidade:  "São Paulo",
		UF:          "SP",
		IBGE:        "3550308",
		GIA:         "1004",
		DDD:         "11",
		SIAFI:       "7107",
	}, nil)

	apiKey := "test_api_key"
	webCEPHandler := web.NewWebCEPHandlerWithDeps(mockCEPRepo, mockWeatherRepo, apiKey)

	validate_cep := usecase.NewValidateCEPUseCase(webCEPHandler.CEPRepository)
	weather_output := validate_cep.Execute(get_weather_dto)
	assert.True(t, weather_output)
}
