package gormdb

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"github.com/snpavlov/gorm-probe/domain"
)

func TestAircraftQuery(t *testing.T) {

	gctx := GormDBContext{}

	dsn := "host=localhost port=5432 dbname=AircraftDemo user=postgres password=Postik123 sslmode=disable"

	db, err := gctx.Open(dsn, "bookings")

	if err != nil {
		t.Fatalf("Cant open context! Error: %v", err)
	}	

	var aircrafts []domain.GAircraft // Declare a slice to hold the results
    result := db.Find(&aircrafts) // Execute the query

	if result.Error != nil {
		t.Fatalf("Cant read data from table! Error: %v", err)
	}

	for _, item := range aircrafts {
		var namelg domain.NameLang

		err := item.JNames.AssignTo(&namelg)
		if err == nil {
    		fmt.Printf("NameRu: '%s', NameEn: '%s'", namelg.NameRu, namelg.NameRu)
		}
    }
}

func TestAirportQuery(t *testing.T) {

	gctx := GormDBContext{}

	dsn := "host=localhost port=5432 dbname=AircraftDemo user=postgres password=Postik123 sslmode=disable"

	db, err := gctx.Open(dsn, "bookings")

	if err != nil {
		t.Fatalf("Cant open context! Error: %v", err)
	}	

	// Create a dry run session
	// dryRunDB := db.Session(&gorm.Session{DryRun: true})

	var airports []domain.GAirport // Declare a slice to hold the results

    result := db.Limit(10).
		Preload("LastDepartureFlights", func(db *gorm.DB) *gorm.DB {
  			return db.Where("actual_departure is not null").
				Order("actual_departure DESC")
		}).
		Preload("LastArrivalFlights", func(db *gorm.DB) *gorm.DB {
  			return db.Where("actual_arrival is not null").
				Order("actual_arrival DESC")
		}).
		Order("airport_code").
		Find(&airports) // Execute the query

	if result.Error != nil {
		t.Fatalf("Cant read data from table! Error: %v", err)
	}

	// Retrieve the generated SQL and variables
	// sql := stmt.SQL.String()
	// vars := stmt.Vars

	// // Print for inspection
	// println("Generated SQL:", sql)
	// println("Variables:", vars)

	for _, item := range airports {
		var aprtname domain.NameLang
		var cityname domain.NameLang

		err := item.JNames.AssignTo(&aprtname)
		if err == nil {
    		fmt.Printf("NameRu: '%s', NameEn: '%s'", aprtname.NameRu, aprtname.NameRu)
		}

		err = item.JCityNames.AssignTo(&cityname)
		if err == nil {
    		fmt.Printf("NameRu: '%s', NameEn: '%s'", cityname.NameRu, cityname.NameRu)
		}
    }
}


// Preload("LastArrivalFlights", func(db *gorm.DB) *gorm.DB {
// 	return db.Joins(`join lateral (
// 			select f.flight_id
// 			from bookings.flights f
// 			where f.arrival_airport = bookings.flights.arrival_airport and f.actual_arrival is not null
// 			order by f.actual_arrival desc
// 			limit 10) as fla on fla.flight_id = bookings.flights.flight_id`)
// }).


func TestFlightsQuery(t *testing.T) {

	gctx := GormDBContext{}

	dsn := "host=localhost port=5432 dbname=AircraftDemo user=postgres password=Postik123 sslmode=disable"

	db, err := gctx.Open(dsn, "bookings")

	if err != nil {
		t.Fatalf("Cant open context! Error: %v", err)
	}	

	var flights []domain.GFlight // Declare a slice to hold the results
    result := db.Limit(10).
		Preload("AirportDeparture").
		Preload("AirportArrival").
		Order("flight_id").
		Find(&flights) // Execute the query

	if result.Error != nil {
		t.Fatalf("Cant read data from table! Error: %v", err)
	}

	for _, item := range flights {
		fmt.Printf("Id: '%v', Code: '%s'", item.Id, item.Code)

		time1 := fmt.Sprintf("Local time:%v", item.PlanDeparture.Time)
		time2 := fmt.Sprintf("Local time:%v", item.PlanDeparture.Time)
		fmt.Printf("t1: %v, t2: %v", time1, time2)
    }
}


// func TestAutomigrate(t *testing.T) {

// 	gctx := GormDBContext{}

// 	dsn := "host=localhost port=5432 dbname=GormTestDB user=RIMDBAdmin password=RimDBAdmin123 sslmode=disable"

// 	err := gctx.Migrate(dsn, "lpok")

// 	if err != nil {
// 		t.Fatalf("Cant migrate database! Error: %v", err)
// 	}

// }
