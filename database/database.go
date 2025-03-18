package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
)

var dir = "./migration"

func Migrate(db *sql.DB) {
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return
		}

		_, err = db.Exec(string(content))

		if err != nil {
			log.Fatal(err)
		}
	}
}

//go:coverage
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")

	return db, err
}

//go:coverage
func ConnectTest() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")

	return db, err
}
