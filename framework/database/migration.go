package database

import (
	"github.com/evolidev/evoli/framework/logging"
	"gorm.io/gorm"
)

type Migration struct {
	callbacks []func(db *gorm.DB)
}

func NewMigrations() *Migration {
	return &Migration{
		callbacks: make([]func(db *gorm.DB), 0),
	}
}

func (m *Migration) Add(migration func(db *gorm.DB)) {
	m.callbacks = append(m.callbacks, migration)
}

func (m *Migration) All() []func(db *gorm.DB) {
	return m.callbacks
}

func (m *Migration) Migrate(db *gorm.DB) {
	for _, migration := range m.callbacks {
		migration(db)

		l := logging.NewLogger(&logging.Config{Name: "db", PrefixColor: 50})
		l.Log("Models migrated successfully")
	}

	m.callbacks = make([]func(db *gorm.DB), 0)
}
