package service

import (
    "testing"
	"log"
    "github.com/snpavlov/app_aircraft/internal/repo"
)

// TestGetAircrafts тестирует получение самолетов
func TestService_GetAircrafts(t *testing.T) {
 
	service, err := AircraftService{}.NewAircraftService(nil)

	if err != nil {
		t.Errorf("Ошибка инициализации сервиса 'AircraftService': %v", err)
	}

	limit := 10

    pager := repo.PageInfo{
        Limit:  &limit,
        Offset: nil, // Без смещения
    }

	result, err := service.GetAircrafts(pager)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAircrafts': %v", err)
    } else 
	{
		log.Printf("Получено %v элементов", len(*result.Data))
	}

}

// TestGetAircrafts тестирует получение самолетов
func TestService_GetAircraftByCodeTest(t *testing.T) {
 
	service, err := AircraftService{}.NewAircraftService(nil)

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

