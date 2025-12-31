package server

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/snpavlov/app_aircraft/internal/auth"
	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/model"
	"github.com/snpavlov/app_aircraft/internal/service"
)


type AppServer struct {
    Greeting *string
	Addr *string
	AppUserKey string
	authService auth.IAppAuthService
	aircraftService service.IAircraftService
	airportService service.IAirportService
}


func (server AppServer) InitConfiguration() (config conf.IConfiguration) {
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("."); 

    if err != nil {
        log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
    }

	return config
}

func (server AppServer) Initialize(usage func()) (AppServer) {

	config := server.InitConfiguration()
	svraddr, err := config.GetServerAddress()
	if err != nil {
		log.Fatalf("Не удалось получить адрес сервера из конфигурации!")
		os.Exit(3)
	}

	server.Greeting = flag.String("g", "Hello", "Greet with `greeting`")
	server.Addr     = flag.String("addr", svraddr, "address to serve")	

	// Parse flags.
	flag.Usage = usage
	flag.Parse()

	// Parse and validate arguments (none).
	args := flag.Args()
	if len(args) != 0 {
		usage()
	}

	server.AppUserKey = "__appuser"

	// Подготка сервиса аутентификации пользователя
	authService, err := auth.AppAuthService{}.NewAppAuthService(config)

	if err != nil {
		log.Fatalf("Ошибка инициализации сервиса 'AppAuthService': %v", err)
		os.Exit(1)
	}

	server.authService = authService

	// Подготка функционального сервиса самолетов
	aircraftService, err := service.AircraftService{}.NewAircraftService(config)

	if err != nil {
		log.Fatalf("Ошибка инициализации сервиса 'AircraftService': %v", err)
		os.Exit(1)
	}

	server.aircraftService = aircraftService

	// Подготка функционального сервиса аэропортов
	airportService, err := service. AirportService{}.NewAirportService(config)

	if err != nil {
		log.Fatalf("Ошибка инициализации сервиса 'AirportService': %v", err)
		os.Exit(1)
	}

	server.airportService = airportService	

	return server

}

func (server AppServer) Greet(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	name := strings.Trim(ctx.Request.URL.Path, "/")
	if name == "" {
		name = "Gopher"
	}

	fmt.Fprintf(ctx.Writer, "<!DOCTYPE html>\n")
	fmt.Fprintf(ctx.Writer, "%s, %s!\n", *server.Greeting, html.EscapeString(name))
}

func (server AppServer) Version(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	info, ok := debug.ReadBuildInfo()
	if !ok {
		ctx.String(500, "no build information available")
		return
	}

	fmt.Fprintf(ctx.Writer, "<!DOCTYPE html>\n<pre>\n")
	fmt.Fprintf(ctx.Writer, "%s\n", html.EscapeString(info.String()))
}

func (server AppServer) Profile(ctx *gin.Context) {

	user, exists := ctx.Get(server.AppUserKey)

	if !exists {
		ctx.IndentedJSON(http.StatusOK, 
			struct { Message string }{Message: "Пользователь не авторизован"})
		return
	}

	// Проверка типа с помощью type assertion
    appuser, ok := user.(auth.AppUser)
    if !ok {
       ctx.IndentedJSON(http.StatusInternalServerError, 
			struct { Message string }{Message: "Данные пользователя не соотвутствуют ожидаемому формату" })
	   return
    }

	if !appuser.IsAuthenticated {
		ctx.IndentedJSON(http.StatusOK, 
			struct { Message string }{Message: "Пользователь не авторизован"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, appuser)
}


func (server AppServer) GetAircafts(ctx *gin.Context) {

	pager := model.PageInfo{
        Limit:  nil,
        Offset: nil,
    }

	err := ctx.ShouldBindQuery(&pager)
	if err != nil {
		argres := model.ServiceListResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка чтения аргументов запроса",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, argres)
		return
	}

	// Call the data method
	result, err := server.aircraftService.GetAircrafts(pager)

	if err != nil {
		result = model.ServiceListResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)
}

func (server AppServer) GetAircaftByCode(ctx *gin.Context) {
	
	code := ctx.Param("code")

	if len(code) == 0 {
		argres := model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка получения шифра. Аргумент 'code' не задан",
		}
		ctx.IndentedJSON(500, argres)
		return
	}	

	// Call the data method
	result, err := server.aircraftService.GetAircraftByCode(code)

	if err != nil {
		result = model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}

func (server AppServer) CreateAircraft(ctx *gin.Context) {
	
	var input model.AircraftInput

	if err := ctx.BindJSON(&input); err != nil {
		argres := model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: fmt.Sprintf("Ошибка получения данных: %v", err.Error()),
		}
		ctx.IndentedJSON(http.StatusBadRequest, argres)
		return
	}

	// Call the data method
	result, err := server.aircraftService.CreateAircraft(input)

	if err != nil {
		result = model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}

func (server AppServer) UpdateAircraft(ctx *gin.Context) {
	
	var input model.AircraftInput

	if err := ctx.BindJSON(&input); err != nil {
		argres := model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: fmt.Sprintf("Ошибка получения данных: %v", err.Error()),
		}
		ctx.IndentedJSON(http.StatusBadRequest, argres)
		return
	}

	// Call the data method
	result, err := server.aircraftService.UpdateAircraft(input)

	if err != nil {
		result = model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}

func (server AppServer) DeleteAircraft(ctx *gin.Context) {
	
	code := ctx.Param("code")

	if len(code) == 0 {
		argres := model.ServiceDataResult[model.AircraftData]{
			Result: false, 
			Message: "Ошибка получения шифра. Аргумент 'code' не задан",
		}
		ctx.IndentedJSON(500, argres)
		return
	}	

	// Call the data method
	result, err := server.aircraftService.DeleteAircraft(code)

	if err != nil {
		result = model.ServiceDataResult[string]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}


func (server AppServer) GetAirports(ctx *gin.Context) {

	pager := model.PageInfo{
        Limit:  nil,
        Offset: nil,
    }

	err := ctx.ShouldBindQuery(&pager)
	if err != nil {
		argres := model.ServiceListResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка чтения аргументов запроса",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, argres)
		return
	}

	// Call the data method
	result, err := server.airportService.GetAirports(pager)

	if err != nil {
		result = model.ServiceListResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)
}

func (server AppServer) GetAirportByCode(ctx *gin.Context) {
	
	code := ctx.Param("code")

	if len(code) == 0 {
		argres := model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка получения шифра. Аргумент 'code' не задан",
		}
		ctx.IndentedJSON(500, argres)
		return
	}	

	// Call the data method
	result, err := server.airportService.GetAirportByCode(code)

	if err != nil {
		result = model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}


func (server AppServer) CreateAirport(ctx *gin.Context) {
	
	var input model.AirportInput

	if err := ctx.BindJSON(&input); err != nil {
		argres := model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: fmt.Sprintf("Ошибка получения данных: %v", err.Error()),
		}
		ctx.IndentedJSON(http.StatusBadRequest, argres)
		return
	}

	// Call the data method
	result, err := server.airportService.CreateAirport(input)

	if err != nil {
		result = model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}

func (server AppServer) UpdateAirport(ctx *gin.Context) {
	
	var input model.AirportInput

	if err := ctx.BindJSON(&input); err != nil {
		argres := model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: fmt.Sprintf("Ошибка получения данных: %v", err.Error()),
		}
		ctx.IndentedJSON(http.StatusBadRequest, argres)
		return
	}

	// Call the data method
	result, err := server.airportService.UpdateAirport(input)

	if err != nil {
		result = model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}

func (server AppServer) DeleteAirport(ctx *gin.Context) {
	
	code := ctx.Param("code")

	if len(code) == 0 {
		argres := model.ServiceDataResult[model.AirportData]{
			Result: false, 
			Message: "Ошибка получения шифра. Аргумент 'code' не задан",
		}
		ctx.IndentedJSON(500, argres)
		return
	}	

	// Call the data method
	result, err := server.airportService.DeleteAirport(code)

	if err != nil {
		result = model.ServiceDataResult[string]{
			Result: false, 
			Message: "Ошибка запроса данных",
			Validations: &[]model.Validation{
				{ Message: fmt.Sprintf("Ошибка: %v", err) },
			},
		}
		ctx.IndentedJSON(500, result)
		return
	}

	ctx.IndentedJSON(http.StatusOK, result)	
}


