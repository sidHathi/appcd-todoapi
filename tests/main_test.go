package tests

import (
	"fmt"
	"os"
	"testing"
	"todo-api/db"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env file for local runs
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using default environment variables")
	}

	// set up database
	db.ConnectForTesting()
	defer db.CleanTestDB()

	// run tests
	code := m.Run()
	os.Exit(code)
}
