package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hawarir/backend-coding-test/controller"
	"github.com/hawarir/backend-coding-test/repository"

	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatalf("Failed to open connection to database: %s", err)
	}
	defer db.Close()

	rideRepo := repository.NewRideRepository(db)

	if err := rideRepo.InitTable(); err != nil {
		log.Fatalf("Failed to initialize table: %s", err)
	}

	e := echo.New()
	controller.SetupRideController(e, rideRepo)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
