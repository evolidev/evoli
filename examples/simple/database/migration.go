package database

import (
	"github.com/evolidev/evoli/examples/simple/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(model.Person{})
}
