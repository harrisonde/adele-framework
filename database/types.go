package database

import "database/sql"

type Database struct {
	DataType string
	Pool     *sql.DB
}

type DataSourceName struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SslMode      string
}
