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
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snpavlov/app_aircraft/internal/server"
	"github.com/snpavlov/app_aircraft/internal/auth"
)

func main() {

	// init app
	server := server.AppServer{}.Initialize(usage)

	// Register handlers.
	router := gin.Default()

	router.Use(server.JwtMiddleware())

	router.GET("/", server.Greet)
	router.GET("/:text", server.Greet)
	router.GET("/profile", server.Profile)
	router.GET("/version", server.Version)

	// Create a group for API version 1
	v1 := router.Group("/api/v1") 
	{
		v1.Use(server.Authorize())

		v1.GET("/aircrafts", server.GetAircafts)
		v1.GET("/aircrafts/:code", server.GetAircaftByCode)

		v1.POST("/aircrafts/create", 
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.CreateAircraft)
		v1.POST("/aircrafts/update", 
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner),
			server.UpdateAircraft)
		v1.POST("/aircrafts/delete/:code",
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.DeleteAircraft)
		v1.DELETE("/aircrafts/:code", 
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner),
			server.DeleteAircraft)

		v1.GET("/airports", server.GetAirports)
		v1.GET("/airports/:code", server.GetAirportByCode)

		v1.POST("/airports/create", 
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.CreateAirport)
		v1.POST("/airports/update", 
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.UpdateAirport)
		v1.POST("/airports/delete/:code",
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.DeleteAirport)
		v1.DELETE("/airports/:code",
			server.Authorize(auth.AppRole_Contrib, auth.AppRole_Owner), 
			server.DeleteAirport)
	}

	startinfo(*server.Addr);

	router.Run(*server.Addr)

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



