package database

import (
	"evoli.dev/framework/config"
	myconfig "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySql struct {
	config *config.Config
}

func (m *MySql) Open() gorm.Dialector {
	dsn := m.buildDSN()

	return mysql.Open(dsn)
}

func (m *MySql) buildDSN() string {
	mysqlConf := myconfig.NewConfig()

	if v := m.config.Get("database").Value().(string); v != "" {
		mysqlConf.DBName = v
	}
	if v := m.config.Get("username").Value().(string); v != "" {
		mysqlConf.User = v
	}
	if v := m.config.Get("password").Value().(string); v != "" {
		mysqlConf.Passwd = v
	}
	if v := m.config.Get("protocol").Value().(string); v != "" {
		mysqlConf.Net = v
	}
	if v := m.config.Get("host").Value().(string); v != "" {
		mysqlConf.Addr = v
	}
	if v := m.config.Get("collation").Value().(string); v != "" {
		mysqlConf.Collation = v
	}

	mysqlConf.Params = make(map[string]string)
	if v := m.config.Get("charset").Value().(string); v != "" {
		mysqlConf.Params["charset"] = v
	}

	mysqlConf.Params["parseTime"] = "True"

	return mysqlConf.FormatDSN()
}
