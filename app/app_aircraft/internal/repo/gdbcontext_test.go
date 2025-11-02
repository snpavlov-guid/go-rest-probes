package repo

import (
	"fmt"
	"testing"
	//"gorm.io/gorm"

	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
	//"github.com/snpavlov/app_aircraft/internal/domain"
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

// TestAirportsQuery
func TestAirportsQuery(t *testing.T) {

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

// TestAirpoertByCodeQuery
func TestAirpoertByCodeQuery(t *testing.T) {
	tests := []struct {
		name    string
		code   string
		expected bool
	}{
		{"TestAircraftByCode_Success", "CNN", true},
		{"TestAircraftByCode_Failed", "XXX", false},
	}

	repo, err := HelperTest_GetAirportRepo() 
	if err != nil {
		t.Errorf("TestAirpoertByCodeQuery: не удалось получить репозиторий: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			airport, err := repo.GetAitportItemByCode(tt.code)
			if err != nil {
				t.Errorf("Ошибка запроса данных 'GetAitportItemByCode': %v", err)
			}

			var rescode string

			if airport == nil {
				rescode = ""
			} else {
				rescode = airport.Code
			}

			restest := rescode == tt.code

			t.Logf("Получен элемент с кодом '%v'", rescode)	

			if restest != tt.expected {
				t.Errorf("Ошибка в тесте %s", tt.name)
			}
		})
	}

}

// TestGetExistsByCode
func TestGetExistsByCodeQuery(t *testing.T) {
	tests := []struct {
		name    string
		code   string
		expected bool
	}{
		{"TestGetExistsByCode_Success", "CNN", true},
		{"TestGetExistsByCode_Failed", "XXX", false},
	}

	repo, err := HelperTest_GetAirportRepo() 
	if err != nil {
		t.Errorf("TestGetExistsByCodeQuery: не удалось получить репозиторий: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			exists, err := repo.GetAitportExistsByCode(tt.code)
			if err != nil {
				t.Errorf("Ошибка запроса данных 'GetAitportExistsByCode': %v", err)
			}

			if exists != tt.expected {
				t.Errorf("Ошибка в тесте %s", tt.name)
			}
		})
	}

}

