package service

import (
	"log"

    "github.com/snpavlov/app_aircraft/internal/conf"
    "github.com/snpavlov/app_aircraft/internal/repo"
	"github.com/snpavlov/app_aircraft/internal/model"
)


// Определяем интерфейс репозитория IAircraftRepo
type IAirportService interface {
	GetAirports(pager model.PageInfo) (model.ServiceListResult[model.AirportData], error)
	GetAirportByCode(code string) (model.ServiceDataResult[model.AirportData], error)
   	// CreateAirport(input model.AircraftInput) (model.ServiceDataResult[model.AirportData], error) 
	// UpdateAirport(input model.AircraftInput) (model.ServiceDataResult[model.AirportData], error) 
	// DeleteAirport(code string) (model.ServiceDataResult[string], error) 
}

type AirportService struct {
    Repo repo.IAirportRepo
}

func (service AirportService) NewAirportService(config conf.IConfiguration) (IAirportService, error) {
       
    // создать экземпляр репозитория
    service.Repo = repo.GormDBContext{Configuration: config};

	return service, nil
}

func (service AirportService) GetAirports(pager model.PageInfo) (model.ServiceListResult[model.AirportData], error) {

	data, total, err := service.Repo.GetAitportItems(pager)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAitportItems': %v", err)
        return model.ServiceListResult[model.AirportData]{}, err
    }

	result := model.ServiceListResult[model.AirportData] { Result: true, Total: total, Items: &data }

	return result, nil
}

func (service AirportService) GetAirportByCode(code string) (model.ServiceDataResult[model.AirportData], error) {

	data, err := service.Repo.GetAitportItemByCode(code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAitportItemByCode': %v", err)
        return model.ServiceDataResult[model.AirportData]{}, err
    }	

	result := model.ServiceDataResult[model.AirportData] { Result: true, Data: data }

	return result, nil
}

