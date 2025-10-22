package repo

import (
	"fmt"
	"testing"
	"gorm.io/gorm"

	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/domain"
	"github.com/snpavlov/app_aircraft/internal/util"

)

func HelperTest_GetAirportRepo() (IAirportRepo, error) {
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

    if err != nil {
        return nil, fmt.Errorf("не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := GormDBContext{Configuration: config};

	return  repo, nil
     
}

// TestAircraftsQuery
func TestAircraftsQuery(t *testing.T) {

	repo, err := HelperTest_GetAirportRepo() 

    if err != nil {
        t.Errorf("Не удалось получить репозиторий: %v", err)
    }

	pager := model.PageInfo{Limit: util.Ptr(10), Offset: util.Ptr(0)}

	airports, total, err := repo.GetAitportItems(pager)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAitportItems': %v", err)
    }

	t.Logf("Получено %v элементов из %v", len(airports), total)

}

// TestAircraftByCodeQuery
func TestAircraftByCodeQuery_Success(t *testing.T) {

	repo, err := HelperTest_GetAirportRepo() 

    if err != nil {
        t.Errorf("Не удалось получить репозиторий: %v", err)
    }

	airportCode := ""

	airport, err := repo.GetAitportItemByCode(airportCode)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAitportItemByCode': %v", err)
    }

	t.Logf("Получен элемент с кодом '%v'", airport.Code)

}

// TestAircraftByCodeQuery
func TestAircraftByCodeQuery_NotFound(t *testing.T) {

	repo, err := HelperTest_GetAirportRepo() 

    if err != nil {
        t.Errorf("Не удалось получить репозиторий: %v", err)
    }

	airportCode := ""

	airport, err := repo.GetAitportItemByCode(airportCode)

	if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAitportItemByCode': %v", err)
    }

	if (airport != nil) {
		t.Errorf("Получен элемент с кодом '%v' всесто nil", airport.Code)
	}

}

