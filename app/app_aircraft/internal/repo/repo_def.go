package repo

import (
	"database/sql"
	"github.com/snpavlov/app_aircraft/internal/model"
)

// Названия типа JSONB
type NameLang struct {
    NameRu string `json:"ru"`
    NameEn string `json:"en"`
}

type Total struct {
	Total  	 int 	 `db:"Total"`
}

type Aircraft struct {
	Code     string  `db:"Code"`
	NameRu   string  `db:"NameRu"`
	NameEn   string  `db:"NameEn"`
	Range  	 int 	 `db:"range"`
}

type SeatType struct {
	Code      string  `db:"Code"`
	SeatType  string  `db:"SeatType"`
	SeatCount int     `db:"SeatCount"`
}

// Определяем интерфейс репозитория IAircraftRepo
type IAircraftRepo interface {
	GetDBConnection() (*sql.DB, error)
	GetAircraftItems(db *sql.DB, pager model.PageInfo) ([]model.AircraftData, int, error)
	GetAircraftItemByCode(db *sql.DB, code string) (*model.AircraftData, error) 
}