package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(filepath string) {
	var err error

	db, err = sql.Open("sqlite3", filepath) // ← здесь убрали :=
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT 0
	)`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create a table : %v", err)
	}
	log.Println("Table checked/created")
}

func InserTask(title string, done bool) (int, error) {
	result, err := db.Exec("INSERT INTO tasks(title, done) VALUES(?,?)", title, done)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
