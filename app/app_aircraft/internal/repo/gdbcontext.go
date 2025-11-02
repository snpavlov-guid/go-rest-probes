package repo

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/jackc/pgx/pgtype"
	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/domain"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/util"
)

var (
	airschema = "bookings"

	airportFlightsQuery = `select * from 
		(select fl.* 
		, 'departure' as source
		, ROW_NUMBER() OVER (PARTITION BY departure_airport ORDER BY actual_departure DESC) AS rownum
		from bookings.flights fl
		where departure_airport in (%s)
		and actual_departure is not null
		union all
		select fl.* 
		, 'arrival' as source
		, ROW_NUMBER() OVER (PARTITION BY arrival_airport ORDER BY actual_arrival DESC) AS rownum
		from bookings.flights fl
		where arrival_airport in (%s)
		and actual_arrival is not null
		) fla
		where fla.rownum <= ?`
)

type GormDBContext struct {
	Configuration conf.IConfiguration
	GormDb *gorm.DB
}

// Определяем интерфейс репозитория IAirportRepo
type IAirportRepo interface {
	GetAitportItems(pager model.PageInfo) ([]model.AirportData, int, error)
	GetAitportItemByCode(code string) (*model.AirportData, error)
	GetAitportExistsByCode(code string) (bool, error)
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

	totalChan := executeGormItemQueryAsync(dctx.GormDb, 
	func(gdb *gorm.DB) (int64, error) {
		var totalCount int64
		result := dctx.GormDb.Model(&domain.GAirport{}).Count(&totalCount)

		return  totalCount, result.Error
	})

	airportsChan := executeGormListQueryAsync(dctx.GormDb, 
		func(gdb *gorm.DB) ([]domain.GAirport, error) {
			var airports []domain.GAirport // Declare a slice to hold the results

			if pager.Offset != nil {
				gdb = gdb.Offset(*pager.Offset)
			}

			if pager.Limit != nil {
				gdb = gdb.Limit(*pager.Limit)
			}

			result := gdb.
				Order("airport_code").
				Find(&airports) // Execute the query

			return airports, result.Error
		})

	airportsRes := <- airportsChan

    if airportsRes.Error != nil {
		return nil, 0, fmt.Errorf("ошибка запроса получения списка Airport: %w", airportsRes.Error)
	}

    // Собрать коды в массив
	codes := util.Map(*airportsRes.Items, func(p domain.GAirport) string {
		return p.Code
	})
	// Сформировать список значений для IN
	inclause := util.GetInClauseItemsString(codes)
	// Сформировать запрос получения полетов из аэропорта
	flightsQuery := fmt.Sprintf(airportFlightsQuery, inclause, inclause)

	airflightChan := executeGormListQueryAsync(dctx.GormDb, 
	func(gdb *gorm.DB) ([]domain.GFlight, error) {
		var airflights []domain.GFlight // Declare a slice to hold the results

		result := gdb.
			Raw(flightsQuery, 10).Find(&airflights) // Execute the query

		return airflights, result.Error
	})

	airflightRes := <- airflightChan

    if airflightRes.Error != nil {
		return nil, 0, fmt.Errorf("ошибка запроса получения списка полетов для Airport: %w", airflightRes.Error)
	}

	fldeparture := util.Filter(*airflightRes.Items, func(fl domain.GFlight) bool {
		return  *fl.Source == "departure"
	})

	flarrival := util.Filter(*airflightRes.Items, func(fl domain.GFlight) bool {
		return  *fl.Source == "arrival"
	})

	totalRes := <- totalChan

	if totalRes.Error != nil {
		return nil, 0, fmt.Errorf("ошибка получения количества записей Airport: %w", airportsRes.Error)
	}

	// Соединяем результаты основного запроса самолетов и данных их мест
    airportItems, err := mapAirportData(*airportsRes.Items, fldeparture, flarrival)
	if err != nil {
		return nil, int(*totalRes.Item),  err
	}

	return airportItems, int(*totalRes.Item), nil
}

func (dctx GormDBContext) GetAitportItemByCode(code string) (*model.AirportData, error) {
	err := dctx.Connect();
    if err != nil {
        return nil, err
    }

	airportChan := executeGormItemQueryAsync(dctx.GormDb, 
		func(gdb *gorm.DB) (domain.GAirport, error) {
			var airport domain.GAirport

			result := dctx.GormDb.
				Where("airport_code = ?", code).
				Assign(domain.GAirport{Code:""}).
				FirstOrInit(&airport) // Execute the query

			return airport, result.Error
		})


	// Сформировать список значений для IN
	inclause := util.GetInClauseItemsString([]string {code})
	// Сформировать запрос получения полетов из аэропорта
	flightsQuery := fmt.Sprintf(airportFlightsQuery, inclause, inclause)

	airflightChan := executeGormListQueryAsync(dctx.GormDb, 
	func(gdb *gorm.DB) ([]domain.GFlight, error) {
		var airflights []domain.GFlight // Declare a slice to hold the results

		result := gdb.
			Raw(flightsQuery, 10).Find(&airflights) // Execute the query

		return airflights, result.Error
	})

	airportRes := <- airportChan

	if airportRes.Error != nil {
		return nil, fmt.Errorf("ошибка запроса получения Airport: %w", err)
	}

	if airportRes.Item.Code == "" {
		return nil, nil
	}	

	airflightRes := <- airflightChan

    if airflightRes.Error != nil {
		return nil, fmt.Errorf("ошибка запроса получения списка полетов для Airport: %w", airflightRes.Error)
	}

	fldeparture := util.Filter(*airflightRes.Items, func(fl domain.GFlight) bool {
		return  *fl.Source == "departure"
	})

	flarrival := util.Filter(*airflightRes.Items, func(fl domain.GFlight) bool {
		return  *fl.Source == "arrival"
	})

	// Соединяем результаты основного запроса самолетов и данных их мест
    airportItem, err := mapAirportItem(*airportRes.Item, fldeparture, flarrival)
	if err != nil {
		return nil, err
	}

	return &airportItem, nil
}

func (dctx GormDBContext) GetAitportExistsByCode(code string) (bool, error) {
	err := dctx.Connect();
    if err != nil {
        return false, err
    }

	var airport domain.GAirport

	result := dctx.GormDb.
			Where("airport_code = ?", code).
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

func mapAirportItem(airport domain.GAirport,
				fldepartures []domain.GFlight, 
				flarrivals []domain.GFlight) (model.AirportData, error) {
	items, err := mapAirportData([]domain.GAirport{airport}, fldepartures, flarrivals)
	if err != nil {
		return model.AirportData{}, err
	}
	return items[0], nil
}

func mapAirportData(airports []domain.GAirport,
				fldepartures []domain.GFlight, 
				flarrivals []domain.GFlight) ([]model.AirportData, error) {
   
    fldeparturesMap := util.SliceToMap(fldepartures, func(p domain.GFlight) string {
		return p.AirportDepartureCode
	})	

    flarrivalsMap := util.SliceToMap(flarrivals, func(p domain.GFlight) string {
		return p.AirportArrivalCode
	})	

    aircraftItems, err := util.Map2(airports, func(p domain.GAirport) (model.AirportData, error) {
		var aprtname domain.NameLang
		var cityname domain.NameLang

		err := p.JNames.AssignTo(&aprtname)
		if err != nil {
			return model.AirportData{}, fmt.Errorf("ошибка парсинга имен аэропорта: %w", err)
		}

		err = p.JCityNames.AssignTo(&cityname)
		if err != nil {
			return model.AirportData{}, fmt.Errorf("ошибка парсинга имен города: %w", err)
		}

		item :=  model.AirportData{ Code: p.Code, 
			NameRu: aprtname.NameRu, 
			NameEn: aprtname.NameEn, 
			CityRu: cityname.NameRu,
			CityEn: cityname.NameEn,
			Timezone: p.Timezone,  } 

		departures, exists := fldeparturesMap[item.Code]
        if exists {
            lastDepartures := util.Map(departures, func(p domain.GFlight) model.AirportFlightData {
                return model.AirportFlightData{Id:p.Id,
					Code: p.Code,
					PlanDeparture: p.PlanDeparture.Time,
					PlanArrival: p.PlanArrival.Time,
					ActualDeparture: util.PrtOrNil(p.ActualDeparture, func(p pgtype.Timestamptz) *time.Time {
						return &p.Time
					}),
					ActualArrival: util.PrtOrNil(p.ActualArrival, func(p pgtype.Timestamptz) *time.Time {
						return &p.Time
					}),
					AircraftCode: p.AircraftCode,
					Status: p.Status,
					AirportDepartureCode: p.AirportDepartureCode,
					AirportArrivalCode: p.AirportArrivalCode,
				}
            })
			item.LastDepartures = &lastDepartures
        }

		arrivals, exists := flarrivalsMap[item.Code]
        if exists {
            lastArrivals := util.Map(arrivals, func(p domain.GFlight) model.AirportFlightData {
                return model.AirportFlightData{Id:p.Id,
					Code: p.Code,
					PlanDeparture: p.PlanDeparture.Time,
					PlanArrival: p.PlanArrival.Time,
					ActualDeparture: util.PrtOrNil(p.ActualDeparture, func(p pgtype.Timestamptz) *time.Time {
						return &p.Time
					}),
					ActualArrival: util.PrtOrNil(p.ActualArrival, func(p pgtype.Timestamptz) *time.Time {
						return &p.Time
					}),					
					AircraftCode: p.AircraftCode,
					Status: p.Status,
					AirportDepartureCode: p.AirportDepartureCode,
					AirportArrivalCode: p.AirportArrivalCode,
				}
            })
			item.LastArrivals = &lastArrivals
        }

        return item, nil
    })

    return aircraftItems, err
}


