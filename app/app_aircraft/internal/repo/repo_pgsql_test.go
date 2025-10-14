package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/util"
)

// TestDBConnection тестирует успешное подключение к базе данных
func TestDBConnection(t *testing.T) {
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

// TestGetAircraftItems тестирует получение комплексного списка самолетов
func TestGetAircraftItems(t *testing.T) {

    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    pager := model.PageInfo{Limit: util.Ptr(5), Offset: util.Ptr(5)}

    aircraftItems, total, err := repo.GetAircraftItems(db, pager)
    if err != nil {
		t.Errorf("Ошибка запроса данных 'GetAircrafts': %v", err)
    }

	t.Logf("Получено %v элементов из %v", len(aircraftItems), total)
	
}

// TestGetAircraftItemByCode тестирует получение комплексного списка самолетов
func TestGetAircraftItemByCode(t *testing.T) {
   
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

	//codeOk := "SU9"
    codeBad := "AN1"


    _, err = repo.GetAircraftItemByCode(db, codeBad)
    if err == nil {
		t.Errorf("Ошибка поиска самолета 'GetAircraft': %v", err)
    } 


}

func TestGetExistsByCodeSuccess(t *testing.T) {
   // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    code := "100"

    exists, err := repo.GetExistsByCode(db, code)
    if err != nil {
		t.Errorf("Ошибка запроса проверки существования самолета 'GetExistsByCode': %v", err)
    } 

   t.Logf("Самолет с кодом '%v': '%v'", code, exists)
}


func TestGetExistsByCodeFail(t *testing.T) {
   // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    code := "AN1"

    exists, err := repo.GetExistsByCode(db, code)
    if err != nil {
		t.Errorf("Ошибка запроса проверки существования самолета 'GetExistsByCode': %v", err)
    } 

   t.Logf("Самолет с кодом '%v': '%v'", code, exists)
}

func TestUpdateAircraft(t *testing.T) {
   // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    input := model.AircraftInput{ Code: "TUS", NameRu: "ТУ 134!", NameEn: "TU 1341", Range: 3531}

    aircraft, err := repo.UpdateAircraft(db, input)
    if err != nil {
		t.Errorf("Ошибка обновления самолета 'UpdateAircraft': %v", err)
    } 

   t.Logf("Обновлен самолет с кодом '%v'", aircraft.Code)
}

func TestDeleteAircraft(t *testing.T) {
   // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    input := "TUS";
    
    code, err := repo.DeleteAircraft(db, input)
    if err != nil {
		t.Errorf("Ошибка удаления самолета 'DeleteAircraft': %v", err)
    } 

   t.Logf("Удален самолет с кодом '%v'", code)
}

func TestCreateAircraft(t *testing.T) {
   // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("./../.."); 

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

    input := model.AircraftInput{ Code: "TUS", NameRu: "ТУ 134", NameEn: "TU 134", Range: 2500}

    aircraft, err := repo.CreateAircraft(db, input)
    if err != nil {
		t.Errorf("Ошибка создания самолета 'CreateAircraft': %v", err)
    } 

   t.Logf("Создан самолет с кодом '%v'", aircraft.Code)
}

