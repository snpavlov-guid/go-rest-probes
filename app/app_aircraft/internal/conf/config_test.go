package conf

import (
	"log"
	"testing"
)

// TestLoadConfiguration тестирует успешное подключение к базе данных
func TestLoadConfiguration(t *testing.T) {
    
    // создать экземпляр конфигурации
    config, err := Configuration{}.New().LoadConfiguration("./../..");

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

    addr, err := config.GetServerAddress()

    if err != nil {
        t.Errorf("Не удалось получить адрес сервера: %v", err)
    }

	log.Printf("Успешное получение конфигурации! Сервер: '%v'", addr)

}

func TestGetPgsqlConnectionString(t *testing.T) {
    
    // создать экземпляр конфигурации
    config, err := Configuration{}.New().LoadConfiguration("./../..");

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

	connection, _ := config.GetPgsqlConnectionString();

	log.Printf("Строка подключения: '%v'", connection)

}
    