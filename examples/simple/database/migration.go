package database

import (
	"gorm.io/gorm"
	"simple/model"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(model.Person{})
}
