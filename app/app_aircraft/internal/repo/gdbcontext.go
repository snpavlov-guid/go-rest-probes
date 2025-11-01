package repo

import(
	"fmt"
	"errors"

	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/snpavlov/app_aircraft/internal/util"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/domain"

)

var (
	airschema = "bookings"
)

type GormDBContext struct {
	Configuration conf.IConfiguration
	GormDb *gorm.DB
}

// Определяем интерфейс репозитория IAirportRepo
type IAirportRepo interface {
	GetAitportItems(pager model.PageInfo) ([]model.AirportData, int, error)
	GetAitportItemByCode(code string) (*model.AirportData, error)
	GetExistsByCode(code string) (bool, error)
	// CreateAircraft(input model.AirportInput) (*model.AirportData, error) 
	// UpdateAircraft(input model.AirportInput) (*model.AirportData, error) 
	// DeleteAircraft(code string) (*string, error) 
}

func (dctx GormDBContext) Open(connection string, dbschema string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{ 
				TablePrefix: dbschema + ".",
			},
	})
	if err != nil {
		return nil, fmt.Errorf("can't open database! Error: %v", err)
	}

	return db, nil
}

func (dctx *GormDBContext) Connect() error {
	
	// Формирование строки подключения
	pgsqlConn, err := dctx.Configuration.GetGormConnectionString()
	if err != nil {
		return fmt.Errorf("ошибка получения строки подключения базы данных: %w", err)
	}

	db, err := dctx.Open(pgsqlConn, airschema)
	
	if err != nil {
		return fmt.Errorf("ошибка подключения к базы данных: %w", err)
	}

	dctx.GormDb = db
	
	return nil
}

func (dctx GormDBContext) GetAitportItems(pager model.PageInfo) ([]model.AirportData, int, error) {
	err := dctx.Connect();
    if err != nil {
        return nil, 0,  err
    }

	var airports []domain.GAirport // Declare a slice to hold the results

    result := dctx.GormDb.
		Offset(*pager.Offset).
		Limit(*pager.Limit).
		Order("airport_code").
		Find(&airports) // Execute the query

	if result.Error != nil {
		return nil, 0, fmt.Errorf("ошибка запроса Airports: %w", err)
	}

	// Соединяем результаты основного запроса самолетов и данных их мест
    airportItems, err := mapAirportData(airports)
	if result.Error != nil {
		return nil, 0,  err
	}

	return airportItems, 0, nil
}

func (dctx GormDBContext) GetAitportItemByCode(code string) (*model.AirportData, error) {
	err := dctx.Connect();
    if err != nil {
        return nil, err
    }

	var airport domain.GAirport

	result := dctx.GormDb.
		Where("code = ?", code).
		Assign(domain.GAirport{Code:""}).
		FirstOrInit(&airport) // Execute the query

	if result.Error != nil {
		return nil, fmt.Errorf("ошибка запроса Airport: %w", err)
	}

	if airport.Code == "" {
		return nil, nil
	}

	// Соединяем результаты основного запроса самолетов и данных их мест
    airportItem, err := mapAirportItem(airport)
	if err != nil {
		return nil, err
	}

	return &airportItem, nil
}

func (dctx GormDBContext) GetExistsByCode(code string) (bool, error) {
	err := dctx.Connect();
    if err != nil {
        return false, err
    }

	var airport domain.GAirport

	result := dctx.GormDb.
			Where("code = ?", code).
			First(&airport) // Execute the query

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Record not found, handle this case (e.g., return a default value or specific error)
			return false, nil
		}
		// Some other database error occurred
		return false, result.Error
	}

	return true, nil

}

func mapAirportItem(airport domain.GAirport) (model.AirportData, error) {
	items, err := mapAirportData([]domain.GAirport{airport})
	if err != nil {
		return model.AirportData{}, err
	}
	return items[0], nil
}

func mapAirportData(airports []domain.GAirport) ([]model.AirportData, error) {
   
    aircraftItems, err := util.Map2(airports, func(p domain.GAirport) (model.AirportData, error) {
		var aprtname domain.NameLang
		var cityname domain.NameLang

		err := p.JNames.AssignTo(&aprtname)
		if err != nil {
			return model.AirportData{}, fmt.Errorf("ошибка парсинга имен аэропорта: %w", err)
		}

		err = p.JCityNames.AssignTo(&cityname)
		if err == nil {
			return model.AirportData{}, fmt.Errorf("ошибка парсинга имен города: %w", err)
		}

		item :=  model.AirportData{ Code: p.Code, 
			NameRu: aprtname.NameRu, 
			NameEn: aprtname.NameEn, 
			CityRu: cityname.NameRu,
			CityEn: cityname.NameEn,
			Timezone: p.Timezone,  } 

        return item, nil
    })

    return aircraftItems, err
}


