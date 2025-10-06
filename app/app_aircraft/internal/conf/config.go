package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

// type Config struct {
// 	dbconnection struct {
// 		host string
// 		port int
// 		username string
// 		password string
// 		database string
// 		sslmode string
// 	}

// 	server struct {
// 		addr string
// 	}
// }

// Определяем интерфейс репозитория IAircraftRepo
type IConfiguration interface {
	LoadConfiguration() (*viper.Viper, error)
	GetPgsqlConnectionString() (string, error)
	GetServerAddress() (string, error)
}

type Configuration struct {
	Rt_viper *viper.Viper
}


func (config Configuration) LoadConfiguration() (*viper.Viper, error) {

    config.Rt_viper.SetConfigName("config") // Имя файла без расширения
    config.Rt_viper.SetConfigType("yml")
    config.Rt_viper.AddConfigPath(".")      // Поиск в корневой директории проекта

    err := config.Rt_viper.ReadInConfig()
    if err != nil {
        return nil, err
    }

    return config.Rt_viper, err
}


func (config Configuration) GetPgsqlConnectionString() (string, error) {

	host := config.Rt_viper.GetString("dbconnection.host")
    port := config.Rt_viper.GetInt("dbconnection.port")
    user := config.Rt_viper.GetString("dbconnection.username")
    password := config.Rt_viper.GetString("dbconnection.password")
    dbname := config.Rt_viper.GetString("dbconnection.database")
    sslmode := config.Rt_viper.GetString("dbconnection.sslmode") // "require" для продакшена

    // Формирование строки подключения
    pgsqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)

    return pgsqlConn, nil
}

func (config Configuration) GetServerAddress() (string, error) {
    svrAddress := config.Rt_viper.GetString("server.addr")
    return svrAddress, nil
}