package model

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Surname     string         `json:"surname"`
	Patronymic  string         `json:"patronymic,omitempty"`
	Gender      string         `json:"gender,omitempty"`
	Age         int            `json:"age,omitempty"`
	Nationality string         `json:"nationality,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
