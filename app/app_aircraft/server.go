// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

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

	"github.com/snpavlov/app_aircraft/internal/conf"
	"github.com/snpavlov/app_aircraft/internal/service"
	"github.com/snpavlov/app_aircraft/internal/model"
)

func main() {

	// init app
	server := AppServer{}.Initialize()

	// Register handlers.
	router := gin.Default()
	router.GET("/", server.greet)
	router.GET("/:text", server.greet)
	router.GET("/version", server.version)

	// Create a group for API version 1
	v1 := router.Group("/api/v1") 
	{
		v1.GET("/aircrafts", server.getAircafts)
		v1.GET("/aircrafts/:code", server.getAircaftByCode)

		v1.POST("/aircrafts/create", server.createAircraft)
		v1.POST("/aircrafts/update", server.updateAircraft)
		v1.POST("/aircrafts/delete/:code", server.deleteAircraft)
		v1.DELETE("/aircrafts/:code", server.deleteAircraft)

		v1.GET("/airports", server.getAirports)
		v1.GET("/airports/:code", server.getAirportByCode)
	}

	startinfo(*server.addr);

	router.Run(*server.addr)

}


type AppServer struct {
    greeting *string
	addr *string
	aircraftService service.IAircraftService
	airportService service.IAirportService
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func startinfo(address string) {
	parts := strings.Split(address, ":")
	if (len(parts[0]) == 0) {
		address = fmt.Sprintf("localhost:%s", parts[1])
	}
	log.Printf("serving http://%s\n", address)
}

func (server AppServer) InitConfiguration() (config conf.IConfiguration) {
    
    // создать экземпляр конфигурации и загрузить данные
    config, err := conf.Configuration{}.New().LoadConfiguration("."); 

    if err != nil {
        log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
    }

	return config
}

func (server AppServer) Initialize() (AppServer) {

	config := server.InitConfiguration()
	svraddr, err := config.GetServerAddress()
	if err != nil {
		log.Fatalf("Не удалось получить адрес сервера из конфигурации!")
		os.Exit(3)
	}

	server.greeting = flag.String("g", "Hello", "Greet with `greeting`")
	server.addr     = flag.String("addr", svraddr, "address to serve")	

	// Parse flags.
	flag.Usage = usage
	flag.Parse()

	// Parse and validate arguments (none).
	args := flag.Args()
	if len(args) != 0 {
		usage()
	}

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

func (server AppServer) greet(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	name := strings.Trim(ctx.Request.URL.Path, "/")
	if name == "" {
		name = "Gopher"
	}

	fmt.Fprintf(ctx.Writer, "<!DOCTYPE html>\n")
	fmt.Fprintf(ctx.Writer, "%s, %s!\n", *server.greeting, html.EscapeString(name))
}

func (server AppServer) version(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	info, ok := debug.ReadBuildInfo()
	if !ok {
		ctx.String(500, "no build information available")
		return
	}

	fmt.Fprintf(ctx.Writer, "<!DOCTYPE html>\n<pre>\n")
	fmt.Fprintf(ctx.Writer, "%s\n", html.EscapeString(info.String()))
}

func (server AppServer) getAircafts(ctx *gin.Context) {

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

func (server AppServer) getAircaftByCode(ctx *gin.Context) {
	
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

func (server AppServer) createAircraft(ctx *gin.Context) {
	
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

func (server AppServer) updateAircraft(ctx *gin.Context) {
	
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

func (server AppServer) deleteAircraft(ctx *gin.Context) {
	
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


func (server AppServer) getAirports(ctx *gin.Context) {

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

func (server AppServer) getAirportByCode(ctx *gin.Context) {
	
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



