package main

import (
	"database/sql"
	"log"

	"github.com/hawarir/backend-coding-test/controller"
	"github.com/hawarir/backend-coding-test/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := sql.Open("sqlite3", "./rides.db")
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

	e.Logger.Fatal(e.Start(":8010"))
}
