package conf

import (
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
    config.rt_viper.SetEnvPrefix("GOAPP")
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
    var dataconn = "dbconnection.data_connection"
    config.rt_viper.BindEnv(dataconn)
    pgsqlConn := config.rt_viper.GetString(dataconn)

    return pgsqlConn, nil
}

func (config Configuration) GetServerAddress() (string, error) {
    var svraddr = "server.addr"
    config.rt_viper.BindEnv(svraddr)
    svrAddress := config.rt_viper.GetString("server.addr")
    return svrAddress, nil
}