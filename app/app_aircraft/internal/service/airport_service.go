package service

import (
	"log"
	"fmt"

    "github.com/snpavlov/app_aircraft/internal/conf"
    "github.com/snpavlov/app_aircraft/internal/repo"
	"github.com/snpavlov/app_aircraft/internal/model"
)


// Определяем интерфейс репозитория IAircraftRepo
type IAirportService interface {
	GetAirports(pager model.PageInfo) (model.ServiceListResult[model.AirportData], error)
	GetAirportByCode(code string) (model.ServiceDataResult[model.AirportData], error)
   	CreateAirport(input model.AirportInput) (model.ServiceDataResult[model.AirportData], error) 
	UpdateAirport(input model.AirportInput) (model.ServiceDataResult[model.AirportData], error) 
	DeleteAirport(code string) (model.ServiceDataResult[string], error) 
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

func (service AirportService) CreateAirport(input model.AirportInput) (model.ServiceDataResult[model.AirportData], error) {
	
    exists, err := service.Repo.GetAitportExistsByCode(input.Code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAitportExistsByCode': %v", err)
        return model.ServiceDataResult[model.AirportData]{}, err
    }

    if (exists) {
        result := model.ServiceDataResult[model.AirportData] { 
            Result: false, 
            Message: fmt.Sprintf("Аэропорт с кодом '%v' уже существует!", input.Code),
         }
        return result, nil
    }

    data, err := service.Repo.CreateAirport(input)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'CreateAirport': %v", err)
        return model.ServiceDataResult[model.AirportData]{}, err
    }

	result := model.ServiceDataResult[model.AirportData] { Result: true, Data: data }

	return result, nil
    
}

func (service AirportService) UpdateAirport(input model.AirportInput) (model.ServiceDataResult[model.AirportData], error) {
	
    exists, err := service.Repo.GetAitportExistsByCode(input.Code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAitportExistsByCode': %v", err)
        return model.ServiceDataResult[model.AirportData]{}, err
    }

    if (!exists) {
        result := model.ServiceDataResult[model.AirportData] { 
            Result: false, 
            Message: fmt.Sprintf("Аэропорт с кодом '%v' не существует!", input.Code),
         }
        return result, nil
    }

    data, err := service.Repo.UpdateAirport(input)
    if err != nil {
		log.Fatalf("Ошибка запроса обновления данных 'UpdateAirport': %v", err)
        return model.ServiceDataResult[model.AirportData]{}, err
    }

	result := model.ServiceDataResult[model.AirportData] { Result: true, Data: data }

	return result, nil
    
}

func (service AirportService) DeleteAirport(code string) (model.ServiceDataResult[string], error) {
	
    exists, err := service.Repo.GetAitportExistsByCode(code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetExistsByCode': %v", err)
        return model.ServiceDataResult[string]{}, err
    }

    if (!exists) {
        result := model.ServiceDataResult[string] { 
            Result: false, 
            Message: fmt.Sprintf("Аэропорт с кодом '%v' не существует!", code),
         }
        return result, nil
    }

    data, err := service.Repo.DeleteAirport(code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'DeleteAirport': %v", err)
        return model.ServiceDataResult[string]{}, err
    }

	result := model.ServiceDataResult[string] { Result: true, Data: data }

	return result, nil
    
}

