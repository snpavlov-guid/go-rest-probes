package conf

import (
	"log"
	"testing"
	"github.com/spf13/viper"
)

// TestLoadConfiguration тестирует успешное подключение к базе данных
func TestLoadConfiguration(t *testing.T) {
    
    // создать экземпляр конфигурации
    config := Configuration{ rt_viper: viper.New()};

	_, err := config.LoadConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

	log.Printf("Успешное получение конфигурации! Сервер: '%v'", 
			config.rt_viper.GetString("server.addr"))

}

func TestGetPgsqlConnectionString(t *testing.T) {
    
    // создать экземпляр конфигурации
    config := Configuration{ rt_viper: viper.New()};

	_, err := config.LoadConfiguration()

    if err != nil {
        t.Errorf("Не удалось загрузить конфигурацию: %v", err)
    }

	connection, _ := config.GetPgsqlConnectionString();

	log.Printf("Строка подключения: '%v'", connection)

}
    