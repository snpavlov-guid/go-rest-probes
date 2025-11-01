package model

import (
	"time"
)

// Данные по классу мест
type SeatData struct {
    SeatType string
    Count int
}

// Общие данные о самолете
type AircraftData struct {
	Code     string  
	NameRu   string 
	NameEn   string  
	Range  	 int 	
	SeatCount int
	Seats *[]SeatData
}

// Общие данные аэропорта
type AirportData struct {
	Code       string 
	NameRu     string
	NameEn     string
	CityRu     string
	CityEn     string	
	Timezone   string

	LastDepartures *[]AirportFlightData
	LastArrivals *[]AirportFlightData
}

// Данные полета аэропорта
type AirportFlightData struct {
	Id int64     
	Code string  
	PlanDeparture time.Time
	PlanArrival   time.Time
	ActualDeparture *time.Time 
	ActualArrival   *time.Time
	AircraftCode string 
	Status string 
	AirportDepartureCode string 
	AirportArrivalCode string
}