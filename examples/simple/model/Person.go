package model

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name string
}
