package repo

import (
	"context"
    "database/sql"
    "time"
    "testing"
    _ "github.com/lib/pq"
    
    "github.com/spf13/viper"
    "github.com/snpavlov/app_aircraft/internal/conf"
)

func CreateConfiguration() (conf.IConfiguration, error) {
    // создать экземпляр конфигурации
    config := conf.Configuration{ Rt_viper: viper.New()}; 

    _, err := config.LoadConfiguration()

     return config, err

}

// TestDBConnection тестирует успешное подключение к базе данных
func TestDBConnection(t *testing.T) {
    
    // создать экземпляр конфигурации
    config, err := CreateConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := AircraftSqlRepo{Configuration: config};
    
    db, err := repo.GetDBConnection()
    if err != nil {
        t.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()
   

    // Проверяем, что подключение действительно работает
    err = db.Ping()
    if err != nil {
        t.Errorf("Ping failed: %v", err)
    }

    // Проверяем, что можем выполнить простой запрос
    var version string
    err = db.QueryRow("SELECT version();").Scan(&version)
    if err != nil {
        t.Errorf("Не удалось выполнить запрос: %v", err)
    }

    t.Logf("Версия PostgreSQL: %s", version)
}

// TestDBConnectionWithTimeout тестирует подключение с таймаутом
func TestDBConnectionWithTimeout(t *testing.T) {
    
    // создать экземпляр конфигурации
    config, err := CreateConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := AircraftSqlRepo{Configuration: config};
     
    db, err := repo.GetDBConnection()
    if err != nil {
        t.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

    // Устанавливаем таймаут для Ping
    db.SetConnMaxLifetime(time.Minute * 3)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = db.PingContext(ctx)
    if err != nil {
        t.Errorf("Ping с таймаутом failed: %v", err)
    }
}

// TestDBConnectionInvalidCredentials тестирует обработку неверных учетных данных
func TestDBConnectionInvalidCredentials(t *testing.T) {
    // Временная функция с неверными данными
    getInvalidConnection := func() (*sql.DB, error) {
        invalidConnStr := "host=localhost port=5432 user=invalid password=wrong dbname=nonexistent sslmode=disable"
        db, err := sql.Open("postgres", invalidConnStr)
        if err != nil {
            return nil, err
        }
        err = db.Ping()
        return db, err
    }

    _, err := getInvalidConnection()
    if err == nil {
        t.Error("Ожидалась ошибка при неверных учетных данных, но ошибки нет")
    }
}

// TestGetAircrafts тестирует получение самолетов
func TestGetAircrafts(t *testing.T) {
 
    // создать экземпляр конфигурации
    config, err := CreateConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := AircraftSqlRepo{Configuration: config};
   
    db, err := repo.GetDBConnection()
    if err != nil {
        t.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	limit := 10

    pager := PageInfo{
        Limit:  &limit,
        Offset: nil, // Без смещения
    }

    withStruct, err := repo.GetAircrafts(db, pager)
    if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAircrafts': %v", err)
    } else 
	{
		t.Logf("Получено %v элементов", len(withStruct))
	}
}

// TestGetAircraft тестирует получение самолета по его коду
func TestGetAircraftSuccess(t *testing.T) {
    
    // создать экземпляр конфигурации
    config, err := CreateConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := AircraftSqlRepo{Configuration: config};
      
    db, err := repo.GetDBConnection()
    if err != nil {
        t.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	codeOk := "SU9"

    _, err = repo.GetAircraftByCode(db, codeOk)
    if err != nil {
		t.Errorf("Ошибка поиска самолета 'GetAircraft': %v", err)
    } 
}

// TestGetAircraft тестирует получение самолета по его коду
func TestGetAircraftFail(t *testing.T) {
 
    // создать экземпляр конфигурации
    config, err := CreateConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    // создать экземпляр репозитория
    repo := AircraftSqlRepo{Configuration: config};
   
    db, err := repo.GetDBConnection()
    if err != nil {
        t.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer db.Close()

	codeBad := "AN1"

    aircraft, err := repo.GetAircraftByCode(db, codeBad)

	if err != nil {
        t.Fatalf("Не удалось выполнить запрос данных: %v", err)
    }

    if aircraft != nil {
        t.Errorf("Ожидалась, что данные не будут получены при неверном коде '%v', но данные есть", codeBad)
    }

}