package service

import (
	"log"
    "fmt"
    "github.com/snpavlov/app_aircraft/internal/conf"
    "github.com/snpavlov/app_aircraft/internal/repo"
	"github.com/snpavlov/app_aircraft/internal/model"
)


// Определяем интерфейс репозитория IAircraftRepo
type IAircraftService interface {
	GetAircrafts(pager model.PageInfo) (model.ServiceListResult[model.AircraftData], error)
	GetAircraftByCode(code string) (model.ServiceDataResult[model.AircraftData], error)
   	CreateAircraft(input model.AircraftInput) (model.ServiceDataResult[model.AircraftData], error) 
	UpdateAircraft(input model.AircraftInput) (model.ServiceDataResult[model.AircraftData], error) 
	DeleteAircraft(code string) (model.ServiceDataResult[string], error) 
}

type AircraftService struct {
    Repo repo.IAircraftRepo
}

func (service AircraftService) NewAircraftService(config conf.IConfiguration) (IAircraftService, error) {
       
    // создать экземпляр репозитория
    service.Repo = repo.AircraftSqlRepo{Configuration: config};

	return service, nil
}

func (service AircraftService) GetAircrafts(pager model.PageInfo) (model.ServiceListResult[model.AircraftData], error) {

	db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
        return model.ServiceListResult[model.AircraftData]{}, err
    }
    defer db.Close()

	data, total, err := service.Repo.GetAircraftItems(db, pager)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircraftItems': %v", err)
        return model.ServiceListResult[model.AircraftData]{}, err
    }

	result := model.ServiceListResult[model.AircraftData] { Result: true, Total: total, Items: &data }

	return result, nil
}

func (service AircraftService) GetAircraftByCode(code string) (model.ServiceDataResult[model.AircraftData], error) {

	db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }
    defer db.Close()

	data, err := service.Repo.GetAircraftItemByCode(db, code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircrafts': %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }

	result := model.ServiceDataResult[model.AircraftData] { Result: true, Data: data }

	return result, nil
}

func (service AircraftService) CreateAircraft(input model.AircraftInput) (model.ServiceDataResult[model.AircraftData], error) {
	
    db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }
    defer db.Close()

    exists, err := service.Repo.GetExistsByCode(db, input.Code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetExistsByCode': %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }

    if (exists) {
        result := model.ServiceDataResult[model.AircraftData] { 
            Result: false, 
            Message: fmt.Sprintf("Самолет с кодом '%v' уже существует!", input.Code),
         }
        return result, nil
    }

    data, err := service.Repo.CreateAircraft(db, input)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'CreateAircraft': %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }

	result := model.ServiceDataResult[model.AircraftData] { Result: true, Data: data }

	return result, nil
    
}

func (service AircraftService) UpdateAircraft(input model.AircraftInput) (model.ServiceDataResult[model.AircraftData], error) {
	
    db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }
    defer db.Close()

    exists, err := service.Repo.GetExistsByCode(db, input.Code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetExistsByCode': %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }

    if (!exists) {
        result := model.ServiceDataResult[model.AircraftData] { 
            Result: false, 
            Message: fmt.Sprintf("Самолет с кодом '%v' не существует!", input.Code),
         }
        return result, nil
    }

    data, err := service.Repo.UpdateAircraft(db, input)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'CreateAircraft': %v", err)
        return model.ServiceDataResult[model.AircraftData]{}, err
    }

	result := model.ServiceDataResult[model.AircraftData] { Result: true, Data: data }

	return result, nil
    
}

func (service AircraftService) DeleteAircraft(code string) (model.ServiceDataResult[string], error) {
	
    db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
        return model.ServiceDataResult[string]{}, err
    }
    defer db.Close()

    exists, err := service.Repo.GetExistsByCode(db, code) 
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetExistsByCode': %v", err)
        return model.ServiceDataResult[string]{}, err
    }

    if (!exists) {
        result := model.ServiceDataResult[string] { 
            Result: false, 
            Message: fmt.Sprintf("Самолет с кодом '%v' не существует!", code),
         }
        return result, nil
    }

    data, err := service.Repo.DeleteAircraft(db, code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'DeleteAircraft': %v", err)
        return model.ServiceDataResult[string]{}, err
    }

	result := model.ServiceDataResult[string] { Result: true, Data: data }

	return result, nil
    
}