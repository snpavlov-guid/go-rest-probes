package repo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/util"
)

var (
	queryAircrafts = `select 
		aircraft_code as "Code"
		, model->>'ru' as "NameRu"
		, model->>'en' as "NameEn"
		, range 
		from bookings.aircrafts_data`
	querySeatTypes = `select 
        aircraft_code as "Code"
        , fare_conditions as "SeatType"
        , count(*) as "SeatCount"
        from bookings.seats st`
	queryTotal = `select count(*) as "Total" from bookings.aircrafts_data`
)

type AircraftSqlRepo struct {
	Configuration conf.IConfiguration
}

// GetDBConnection возвращает подключение к PostgreSQL и ошибку.
// Параметры подключения лучше вынести в конфигурацию (например, через флаги или переменные окружения).
func (repo AircraftSqlRepo) GetDBConnection() (*sql.DB, error) {

	// Формирование строки подключения
	pgsqlConn, err := repo.Configuration.GetPgsqlConnectionString()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения строки подключения базы данных: %w", err)
	}

	// Открытие подключения (пула соединений)
	db, err := sql.Open("postgres", pgsqlConn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	// Опционально: конфигурация пула соединений
	//db.SetMaxOpenConns(25) // Максимум открытых соединений
	//db.SetMaxIdleConns(25) // Максимум соединений в режиме ожидания
	//db.SetConnMaxLifetime(0) // Время жизни соединения (0 = бессрочно)

	log.Println("Успешное подключение к базе данных!")
	return db, nil
}

// GetAircraftItems возвращает самолеты с пагинацией
func (repo AircraftSqlRepo) GetAircraftItems(db *sql.DB, pager model.PageInfo) ([]model.AircraftData, int, error) {

    query := util.AddOrderByClause(queryAircrafts, []model.OrderInfo{{Field: "Code"}})
	query, args := util.AddPaginationClause(query, pager)

    aircrafts, err := executeRowsQuery(db, query, args, 
        func(rows *sql.Rows) (Aircraft, error) {
            var item Aircraft
			err := rows.Scan(
			    &item.Code,
			    &item.NameRu,
			    &item.NameEn,
			    &item.Range,
            )
			return item, err
		},
    )

    if err != nil {
		return nil, 0, fmt.Errorf("ошибка запроса Aircraft: %w", err)
	}

    // Собрать коды в массив
	codes := util.Map(aircrafts, func(p Aircraft) string {
		return p.Code
	})

    // Готовим запрос на места
	query = util.AddInClause(querySeatTypes, codes, "aircraft_code", "WHERE")
    query = util.AddGroupClause(query, []string{"aircraft_code", "fare_conditions"})
	query = util.AddOrderByClause(query, []model.OrderInfo{{Field: "Code"}, {Field: "SeatType"}})

	var arg0 []any
    seatTypes, err := executeRowsQuery(db, query, arg0, 
        func(rows *sql.Rows) (SeatType, error) {
            var item SeatType
			err := rows.Scan(
			    &item.Code,
			    &item.SeatType,
			    &item.SeatCount,               
            )
			return item, err
		},
    )

    if err != nil {
		return nil, 0, fmt.Errorf("ошибка запроса SeatType: %w", err)
	}

    // Соединяем результаты основного запроса самолетов и данных их мест
    aircraftItems := mapAircraftData(aircrafts, seatTypes)

    total, err := executeRowQuery(db, queryTotal, arg0, 
        func(row *sql.Row) (Total, error) {
            var item Total
			err := row.Scan(
			    &item.Total,
            )
			return item, err
		},
    )

    if err != nil {
		return nil, 0, fmt.Errorf("ошибка запроса Total: %w", err)
	}  

	return aircraftItems, total.Total, nil
}

// GetAircraftItems возвращает самолеты с пагинацией
func (repo AircraftSqlRepo) GetAircraftItemByCode(db *sql.DB, code string) (*model.AircraftData, error) {

	query := util.AddWhereClause(queryAircrafts, []string{"aircraft_code"}, 1, "WHERE", "AND")

	args := []any{code}
    aircraft, err := executeRowQuery(db, query, args, 
        func(row *sql.Row) (Aircraft, error) {
			var item Aircraft
			err := row.Scan(
			    &item.Code,
			    &item.NameRu,
			    &item.NameEn,
			    &item.Range,
            )
			return item, err
		},
    )

    if err != nil {
		return nil, fmt.Errorf("ошибка запроса Aircraft: %w", err)
	}  	

	if aircraft == nil {
		return nil, nil
	}


    // Готовим запрос на места
	query = util.AddInClause(querySeatTypes, []string{aircraft.Code}, "aircraft_code", "WHERE")
    query = util.AddGroupClause(query, []string{"aircraft_code", "fare_conditions"})
	query = util.AddOrderByClause(query, []model.OrderInfo{{Field: "Code"}, {Field: "SeatType"}})

	var arg0 []any
    seatTypes, err := executeRowsQuery(db, query, arg0, 
        func(rows *sql.Rows) (SeatType, error) {
            var item SeatType
			err := rows.Scan(
			    &item.Code,
			    &item.SeatType,
			    &item.SeatCount,               
            )
			return item, err
		},
    )

    if err != nil {
		return nil, fmt.Errorf("ошибка запроса SeatType: %w", err)
	}

	// Соединяем результаты основного запроса самолетов и данных их мест
	aircraftItem := mapAircraftItem(*aircraft, seatTypes);

	return &aircraftItem, nil
}


func mapAircraftItem(aircraft Aircraft, seatTypes []SeatType) (model.AircraftData) {
	items := mapAircraftData([]Aircraft{aircraft}, seatTypes)
	return items[0]
}

func mapAircraftData(aircrafts []Aircraft, seatTypes []SeatType) ([]model.AircraftData) {
	
    seatMap := util.SliceToMap(seatTypes, func(p SeatType) string {
		return p.Code
	})	
    
    aircraftItems := util.Map(aircrafts, func(p Aircraft) model.AircraftData {
		item :=  model.AircraftData{ Code: p.Code, NameRu: p.NameRu, NameEn: p.NameEn, Range: p.Range } 
        seats, exists := seatMap[item.Code]
        if exists {
            seatItems := util.Map(seats, func(p SeatType) model.SeatData {
                return model.SeatData{SeatType: p.SeatType, Count: p.SeatCount}
            })
            item.Seats = &seatItems
            item.SeatCount = util.Sum(seatItems, func(p model.SeatData) int {
                return p.Count
            })
        }
        return item
    })

    return aircraftItems
}

