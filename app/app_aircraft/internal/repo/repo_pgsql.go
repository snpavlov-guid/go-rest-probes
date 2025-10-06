package repo

import (
	"database/sql"
	"fmt"
	"log"
    _ "github.com/lib/pq"
    "github.com/snpavlov/app_aircraft/internal/conf"

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

// GetAircraftsWithPagination возвращает самолеты с пагинацией
func (repo AircraftSqlRepo) GetAircraftsWithPagination(db *sql.DB, limit, offset int) ([]Aircraft, error) {
    query := `select 
		aircraft_code as "Code"
		, model->'ru' as "NameRu"
		, model->'en' as "NameEn"
		, range 
		from bookings.aircrafts_data ORDER BY aircraft_code 
		LIMIT $1 OFFSET $2`
    
    rows, err := db.Query(query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
    }
    defer rows.Close()
    
    var aircrafts []Aircraft
    for rows.Next() {
        var aircraft Aircraft
        if err := rows.Scan(
            &aircraft.Code,
            &aircraft.NameRu,
            &aircraft.NameEn,
            &aircraft.Range,
        ); err != nil {
            return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
        }
        aircrafts = append(aircrafts, aircraft)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
    }
    
    return aircrafts, nil
}

// GetAircraftsWithPagination возвращает самолеты с пагинацией
func (repo AircraftSqlRepo) GetAircrafts(db *sql.DB, pager PageInfo) ([]Aircraft, error) {
    query := `select 
		aircraft_code as "Code"
		, model->'ru' as "NameRu"
		, model->'en' as "NameEn"
		, range 
		from bookings.aircrafts_data ORDER BY aircraft_code`
    
	var args []interface{}
    paramCount := 0
    
    if pager.Limit != nil {
        paramCount++
        query += fmt.Sprintf(" LIMIT $%d", paramCount)
        args = append(args, *pager.Limit)
    }
    
    if pager.Offset != nil {
        paramCount++
        query += fmt.Sprintf(" OFFSET $%d", paramCount)
        args = append(args, *pager.Offset)
    }
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
    }
    defer rows.Close()
    
    var aircrafts []Aircraft
    for rows.Next() {
        var aircraft Aircraft
        if err := rows.Scan(
            &aircraft.Code,
            &aircraft.NameRu,
            &aircraft.NameEn,
            &aircraft.Range,
        ); err != nil {
            return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
        }
        aircrafts = append(aircrafts, aircraft)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
    }
    
    return aircrafts, nil
}

// GetAircraftByCode возвращает самолет по коду
func (repo AircraftSqlRepo) GetAircraftByCode(db *sql.DB, code string) (*Aircraft, error) {
    query := `select 
		aircraft_code as "Code"
		, model->'ru' as "NameRu"
		, model->'en' as "NameEn"
		, range 
		from bookings.aircrafts_data WHERE aircraft_code = $1`
    
    var aircraft Aircraft
    err := db.QueryRow(query, code).Scan(
        &aircraft.Code,
        &aircraft.NameRu,
        &aircraft.NameEn,
        &aircraft.Range,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Самолет не найден
        }
        return nil, err
    }
    
    return &aircraft, nil
}


