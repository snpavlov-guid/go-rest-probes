package repo

import (
	"database/sql"
)


type PageInfo struct {
    Limit  *int `form:"size"`
    Offset *int `form:"offset"`
}

// Названия типа JSONB
type NameLang struct {
    NameRu string `json:"ru"`
    NameEn string `json:"en"`
}

type Aircraft struct {
	Code     string  `db:"Code"`
	NameRu   string  `db:"NameRu"`
	NameEn   string  `db:"NameEn"`
	Range  	 int 	 `db:"range"`
}


// Определяем интерфейс репозитория IAircraftRepo
type IAircraftRepo interface {
	GetDBConnection() (*sql.DB, error)
	GetAircrafts(db *sql.DB, pager PageInfo) ([]Aircraft, error)
	GetAircraftByCode(db *sql.DB, code string) (*Aircraft, error)
}