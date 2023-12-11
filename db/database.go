package db

import (
	"database/sql"

	"github.com/bryansamperura/ticket-booking/config"
	"github.com/bryansamperura/ticket-booking/helper"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Init() {
	conf := config.GetConfig()

	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString)

	helper.PanicIfError(err)

	err = db.Ping()

	helper.PanicIfError(err)
}

func CreateConnection() *sql.DB {
	return db
}
