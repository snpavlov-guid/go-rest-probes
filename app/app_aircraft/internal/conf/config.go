package conf

import (
	"fmt"
	"github.com/spf13/viper"
)


// Определяем интерфейс репозитория IAircraftRepo
type IConfiguration interface {
	LoadConfiguration(basePath string) (IConfiguration, error)
	GetPgsqlConnectionString() (string, error)
	GetServerAddress() (string, error)
}

type Configuration struct {
	rt_viper *viper.Viper
}

func (config Configuration) New() (IConfiguration) {
    config.rt_viper = viper.New()
    return config;
}

func (config Configuration) LoadConfiguration(basePath string) (IConfiguration, error) {

    config.rt_viper.SetConfigName("config") // Имя файла без расширения
    config.rt_viper.SetConfigType("yml")
    config.rt_viper.AddConfigPath(basePath) // Директория конфигурации

    err := config.rt_viper.ReadInConfig()
    if err != nil {
        return nil, err
    }

    return config, err
}


func (config Configuration) GetPgsqlConnectionString() (string, error) {

	host := config.rt_viper.GetString("dbconnection.host")
    port := config.rt_viper.GetInt("dbconnection.port")
    user := config.rt_viper.GetString("dbconnection.username")
    password := config.rt_viper.GetString("dbconnection.password")
    dbname := config.rt_viper.GetString("dbconnection.database")
    sslmode := config.rt_viper.GetString("dbconnection.sslmode") // "require" для продакшена

    // Формирование строки подключения
    pgsqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)

    return pgsqlConn, nil
}

func (config Configuration) GetServerAddress() (string, error) {
    svrAddress := config.rt_viper.GetString("server.addr")
    return svrAddress, nil
}