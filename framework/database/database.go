package database

import (
	"github.com/evolidev/evoli/framework/config"
	"gorm.io/gorm"
)

type Driver interface {
	Open() gorm.Dialector
}

type Database struct {
	connection Driver
	config     *config.Config
}

func (d *Database) Connect() *gorm.DB {
	db, _ := gorm.Open(d.connection.Open(), &gorm.Config{})

	return db
}

func Get(config *config.Config) *Database {
	d := driver(config)

	return &Database{connection: d, config: config}
}

func driver(config *config.Config) Driver {
	switch config.Get("default").Value().(string) {
	case "sqlite":
		myConf := config.Get("sqlite")
		myConf.Set("base", config.Get("base").Value().(string))

		return &Sqlite{
			config: myConf,
		}
	case "mysql":
		return &MySql{config: config.Get("mysql")}
	}

	panic("unknown database configured")
}
