package domain

import (
	"time"
)

// Пьеса
type Play struct {
	Id 			int64   `gorm:"primaryKey;autoIncrement"`
	Name     	string  `gorm:"size:512;index:IDX_PlayName;not null"`
	Description string  `gorm:"size:1536"`
	Audience 	int     `gorm:"index:IDX_PlayAudience"`
	Actors		[]Actor `gorm:"many2many:PlayActor;"`
}

// Актер 
type Actor struct {
	Id 	     int64   `gorm:"primaryKey;autoIncrement"`
	Name     string  `gorm:"size:64;index:IDX_ActorName;not null"`
	Surname  string  `gorm:"size:64;index:IDX_ActorSurname;not null"`
	Bithdate time.Time
}

// Показ
type Showing struct {
	Id 	    int64 	  `gorm:"primaryKey;autoIncrement"`
	PlayId  int64
	Play    Play	  `gorm:"foreignkey:PlayId"`
	Date    time.Time `gorm:"index:IDX_ShowingDate;not null"`
	Status  byte      `gorm:"index:IDX_ShowingStatus"`
	Address string    `gorm:"size:1024;index:IDX_ShowingAddress"`
	Mark    float32   `gorm:"type:numeric(5,2);index:IDX_ShowingMark"`
}