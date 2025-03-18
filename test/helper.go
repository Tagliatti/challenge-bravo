package test

import (
	"database/sql"
	"github.com/Tagliatti/challenge-bravo/database"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var currentDir, _ = os.Getwd()
var baseDir = filepath.Join(currentDir, "../..")

type ApiTestCase struct {
	Name         string
	Content      string
	ExpectedCode int
	ExpectedBody string
}

func SetUp(t *testing.T) *sql.DB {
	db, err := SetUpTestDatabase(baseDir)

	if err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		if db != nil {
			db.Close()
		}
	})

	return db
}

func SetUpTestDatabase(projectDir string) (*sql.DB, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if currentDir != projectDir {
		os.Chdir(projectDir)
	}

	db, err := database.ConnectTest()

	if err != nil {
		return nil, err
	}

	database.Migrate(db)

	return db, nil
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewHttpClientWithError(t *testing.T) *http.Client {
	return &http.Client{Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		t.Error("currency api handler should not be called")
		return nil, nil
	})}
}
