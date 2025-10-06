package model

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