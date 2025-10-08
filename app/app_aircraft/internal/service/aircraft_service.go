package service

import (
	"log"
    "github.com/snpavlov/app_aircraft/internal/conf"
    "github.com/snpavlov/app_aircraft/internal/repo"
	"github.com/snpavlov/app_aircraft/internal/model"
)


// Определяем интерфейс репозитория IAircraftRepo
type IAircraftService interface {
	GetAircrafts(pager model.PageInfo) (model.ServiceListResult[model.AircraftData], error)
	GetAircraftByCode(code string) (model.ServiceDataResult[model.AircraftData], error)
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
    }
    defer db.Close()

	data, total, err := service.Repo.GetAircraftItems(db, pager)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircraftItems': %v", err)
    }

	result := model.ServiceListResult[model.AircraftData] { Result: true, Total: total, Items: &data }

	return result, nil
}

func (service AircraftService) GetAircraftByCode(code string) (model.ServiceDataResult[model.AircraftData], error) {

	db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	data, err := service.Repo.GetAircraftItemByCode(db, code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircrafts': %v", err)
    }

	result := model.ServiceDataResult[model.AircraftData] { Result: true, Data: data }

	return result, nil
}