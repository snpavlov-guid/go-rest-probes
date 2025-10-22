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
    Code     string  `json:"code"`
    SeatType string  `json:"seatType"`
    SeatNumb string  `json:"seatNumb"`

}

// Данные названия
type NameInput struct {
	En   string  `json:"en"`
	Ru   string  `json:"ru"`
}

// Общие данные о самолете
type AircraftInput struct {
	Code     string  `json:"code"`
	NameRu   string  `json:"nameRu"`
	NameEn   string  `json:"nameEn"`
	Range  	 int 	 `json:"range"`
}

// Общие данные аэропорта
type AirportInput struct {
	Code       string 
	NameRu     string
	NameEn     string
	CityRu     string
	CityEn     string	
	Timezone   string
}