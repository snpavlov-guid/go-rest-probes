package domain

import (
	//"time"
	"github.com/jackc/pgx/pgtype"
) 

// Названия типа JSONB
type NameLang struct {
    NameRu string `json:"ru"`
    NameEn string `json:"en"`
}


type GAircraft struct {
	Code     string  `gorm:"primaryKey;column:aircraft_code;not null"`
	JNames   pgtype.JSONB `gorm:"type:jsonb;default:'{}';column:model;not null"`
	Range  	 int 	 `gorm:"column:range;not null"`
}

// TableName specifies the table name for the User model
func (GAircraft) TableName() string {
	return "aircrafts_data" // Your desired table name
}

type GAirport struct {
	Code       string  `gorm:"primaryKey;column:airport_code;not null"`
	JNames     pgtype.JSONB `gorm:"type:jsonb;default:'{}';column:airport_name;not null"`
	JCityNames pgtype.JSONB `gorm:"type:jsonb;default:'{}';column:city;not null"`
	Position   Point `gorm:"type:point;column:coordinates;not null"` 
	Timezone   string `gorm:"column:timezone;not null"`

	LastDepartureFlights *[]GFlight `gorm:"foreignKey:AirportDepartureCode"`
	LastArrivalFlights *[]GFlight `gorm:"foreignKey:AirportArrivalCode"`
}

// TableName specifies the table name for the User model
func (GAirport) TableName() string {
	return "airports_data" // Your desired table name
}

type GFlight struct {
	Id int64     `gorm:"primaryKey;column:flight_id;not null"`
	Code string  `gorm:"column:flight_no;not null"`
	PlanDeparture pgtype.Timestamptz `gorm:"column:scheduled_departure;not null"`
	PlanArrival pgtype.Timestamptz   `gorm:"column:scheduled_arrival;not null"`
	ActualDeparture* pgtype.Timestamptz `gorm:"column:actual_departure"`
	ActualArrival* pgtype.Timestamptz   `gorm:"column:actual_arrival"`
	AircraftCode string `gorm:"column:aircraft_code;not null"`
	Status string `gorm:"column:status;not null"`
	AirportDepartureCode string `gorm:"column:departure_airport;not null"`
	AirportArrivalCode string `gorm:"column:arrival_airport;not null"`
	AirportDeparture* GAirport `gorm:"foreignkey:AirportDepartureCode"`
	AirportArrival* GAirport `gorm:"foreignkey:AirportArrivalCode"`
}

// TableName specifies the table name for the User model
func (GFlight) TableName() string {
	return "flights" // Your desired table name
}
