package use

import (
	"evoli.dev/framework/database"
	"gorm.io/gorm"
	"time"
)

var myDb *gorm.DB
var migration *database.Migration

func Migration() *database.Migration {
	if nil == migration {
		migration = database.NewMigrations()
	}

	return migration
}

func DB() *gorm.DB {
	if nil != myDb {
		return myDb
	}

	dbConfig := Config("db")
	dbConfig.Set("base", BasePath())
	myDb = database.Get(dbConfig).Connect()

	configurePooling()

	return myDb
}

func configurePooling() {
	sqlDB, _ := myDb.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
