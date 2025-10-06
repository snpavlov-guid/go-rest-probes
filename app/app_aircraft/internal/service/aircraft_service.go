package service

import (
	"log"
    "github.com/spf13/viper"
    "github.com/snpavlov/app_aircraft/internal/conf"
    "github.com/snpavlov/app_aircraft/internal/repo"
)

// Названия типа JSONB
type Validation struct {
    Property string
    Message string
}

type ServiceDataResult[TD any] struct {
	Result   bool
	Message  string
	Validations* []Validation
	Code *string
	Data *TD
}

// Определяем интерфейс репозитория IAircraftRepo
type IAircraftService interface {
	GetAircrafts(pager repo.PageInfo) (ServiceDataResult[[]repo.Aircraft], error)
	GetAircraftByCode(code string) (ServiceDataResult[*repo.Aircraft], error)
}

type AircraftService struct {
    Repo repo.IAircraftRepo
}

func (service AircraftService) NewAircraftService(extconfig *conf.Configuration) (IAircraftService, error) {
    // создать экземпляр конфигурации
    config := conf.Configuration{ Rt_viper: viper.New()}; 

    if extconfig != nil {
        config = *extconfig
    }

	_, err := config.LoadConfiguration()
	
    if err != nil {
        log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    service.Repo = repo.AircraftSqlRepo{Configuration: config};

	return service, nil
}

func (service AircraftService) GetAircrafts(pager repo.PageInfo) (ServiceDataResult[[]repo.Aircraft], error) {

	db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	data, err := service.Repo.GetAircrafts(db, pager)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircrafts': %v", err)
    }

	result := ServiceDataResult[[]repo.Aircraft] { Result: true, Data: &data }

	return result, nil
}

func (service AircraftService) GetAircraftByCode(code string) (ServiceDataResult[*repo.Aircraft], error) {

	db, err := service.Repo.GetDBConnection()
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	data, err := service.Repo.GetAircraftByCode(db, code)
    if err != nil {
		log.Fatalf("Ошибка запроса данных 'GetAircrafts': %v", err)
    }

	result := ServiceDataResult[*repo.Aircraft] { Result: true, Data: &data }

	return result, nil
}