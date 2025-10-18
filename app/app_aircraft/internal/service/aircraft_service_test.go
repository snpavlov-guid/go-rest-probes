package service

import (
    "testing"
	"log"
	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
)

// TestGetAircrafts тестирует получение самолетов
func TestService_GetAircrafts(t *testing.T) {

    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }
	
	service, err := AircraftService{}.NewAircraftService(config)

	if err != nil {
		t.Errorf("Ошибка инициализации сервиса 'AircraftService': %v", err)
	}

	limit := 10

    pager := model.PageInfo{
        Limit:  &limit,
        Offset: nil, // Без смещения
    }

	result, err := service.GetAircrafts(pager)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAircrafts': %v", err)
    } else 
	{
		log.Printf("Получено %v элементов", len(*result.Items))
	}

}

// TestGetAircrafts тестирует получение самолетов
func TestService_GetAircraftByCodeTest(t *testing.T) {
 
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

	service, err := AircraftService{}.NewAircraftService(config)

	if err != nil {
		t.Errorf("Ошибка инициализации сервиса 'AircraftService': %v", err)
	}

	codeOk := "SU9"

	result, err := service.GetAircraftByCode(codeOk)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAircrafts': %v", err)
    } else 
	{
		log.Printf("Получен объект '%v'", *result.Data)
	}

}

