package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Database struct {
	Db          *sql.DB
	OrdersTable OrdersTable
}

func NewDatabase() *Database {
	db := newDB()
	return &Database{db, newOrdersTable(db)}
}

type NotFoundError struct {
}

func (err NotFoundError) Error() string {
	return "not found"
}

func newDB() *sql.DB {
	viper.SetDefault("sql.hostname", "localhost")
	viper.SetDefault("sql.port", "5432")
	viper.SetDefault("sql.user", "postgres")
	viper.SetDefault("sql.sslmode", "disable")

	host := viper.GetString("sql.hostname")
	port := viper.GetString("sql.port")
	user := viper.GetString("sql.user")
	password := viper.GetString("sql.password")
	databaseName := viper.GetString("sql.database")
	sslmode := viper.GetString("sql.sslmode")

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s"+
		" sslmode=%s", host, port, user, password, databaseName, sslmode)

	db, err := sql.Open("postgres", info)
	if err != nil {
		log.Fatalln("Error occurred during opening database connection: ",
			err.Error())
	}

	return db
}
