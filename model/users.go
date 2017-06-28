package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Uuid             string `gorm:"primary_key"`
	Login            string `gorm:"primary_key"`
	RegistrationDate time.Time
}