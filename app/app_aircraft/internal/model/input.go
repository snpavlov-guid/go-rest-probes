package model

type PageInfo struct {
    Limit  *int `form:"size"`
    Offset *int `form:"offset"`
}

type OrderInfo struct {
    Field  string
    Desc bool
}


// Данные по классу мест
type SeatInput struct {
    Code     string  
    SeatType string
    SeatNumb string

}

// Данные названия
type NameInput struct {
	En   string  `json:"en"`
	Ru   string  `json:"ru"`
}

// Общие данные о самолете
type AircraftInput struct {
	Code     string  
	NameRu   string 
	NameEn   string  
	Range  	 int 	
}