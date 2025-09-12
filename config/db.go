package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=localhost port=5432 user=%s password=%s dbname=todo_app sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	createTable :=
		`CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		completed BOOLEAN NOT NULL,
		category VARCHAR(50) NOT NULL,
		priority VARCHAR(10) NOT NULL,
		completed_at TIMESTAMP NULL,
		due_date TIMESTAMP NULL
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}
}
