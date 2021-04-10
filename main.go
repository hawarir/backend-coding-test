package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

var (
	tableSchema     map[string]string
	tableColumns    []string
	tableDefinition []string
)

func init() {
	tableSchema = map[string]string{
		"id":            "INTEGER PRIMARY KEY AUTOINCREMENT",
		"startLat":      "REAL",
		"startLong":     "REAL",
		"endLat":        "REAL",
		"endLong":       "REAL",
		"riderName":     "TEXT",
		"driverName":    "TEXT",
		"driverVehicle": "TEXT",
	}

	tableColumns = make([]string, len(tableSchema))
	tableDefinition = make([]string, len(tableSchema))

	i := 0
	for key, value := range tableSchema {
		tableColumns[i] = key
		tableDefinition[i] = fmt.Sprintf("%s %s", key, value)
		i++
	}
}

type (
	Ride struct {
		ID             int64   `json:"id"`
		StartLatitude  float64 `json:"startLatitude"`
		StartLongitude float64 `json:"startLongitude"`
		EndLatitude    float64 `json:"endLatitude"`
		EndLongitude   float64 `json:"endLongitude"`
		RiderName      string  `json:"riderName"`
		DriverName     string  `json:"driverName"`
		DriverVehicle  string  `json:"driverVehicle"`
	}

	Handler struct {
		db *sql.DB
	}
)

func (r Ride) validate() error {
	errs := []string{}
	correctLatitude := func(lat float64) bool {
		return lat >= -90 && lat <= 90
	}
	correctLongitude := func(long float64) bool {
		return long >= -180 && long <= 180
	}
	stringEmpty := func(s string) bool {
		return s == ""
	}
	for _, lat := range []float64{r.StartLatitude, r.EndLatitude} {
		if !correctLatitude(lat) {
			errs = append(errs, fmt.Sprintf("%f is not a valid latitude value", lat))
		}
	}
	for _, long := range []float64{r.StartLongitude, r.EndLongitude} {
		if !correctLongitude(long) {
			errs = append(errs, fmt.Sprintf("%f is not a valid longitude value", long))
		}
	}
	for key, value := range map[string]string{
		"riderName":     r.RiderName,
		"driverName":    r.DriverName,
		"driverVehicle": r.DriverVehicle,
	} {
		if stringEmpty(value) {
			errs = append(errs, fmt.Sprintf("%s can't be empty", key))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func newHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func main() {
	db, err := sql.Open("sqlite3", "./rides.db")
	if err != nil {
		log.Fatalf("Failed to open connection to database: %s", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS rides (" + strings.Join(tableDefinition, ",") + ")"); err != nil {
		log.Fatalf("Failed to initialize table: %s", err)
	}

	h := newHandler(db)

	e := echo.New()
	e.GET("/health", healthCheck)

	e.POST("/rides", h.addRide)
	e.GET("/rides", h.getAllRides)
	e.GET("/rides/:id", h.getRide)

	e.Logger.Fatal(e.Start(":8010"))
}

func healthCheck(c echo.Context) error {
	return c.String(200, "Healthy")
}

func (h *Handler) addRide(c echo.Context) error {
	var ride Ride
	if err := c.Bind(&ride); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Malformed request body: %s", err))
	}
	if err := ride.validate(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid request body: %s", err))
	}
	lastInsertID, err := insertRide(h.db, ride)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	ride.ID = lastInsertID
	return c.JSON(http.StatusCreated, ride)
}

func (h *Handler) getAllRides(c echo.Context) error {
	rides, err := selectAllRides(h.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	return c.JSON(http.StatusOK, rides)
}

func (h *Handler) getRide(c echo.Context) error {
	id := c.Param("id")
	rideID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid ID: %s", err))
	}
	ride, err := selectRideByID(h.db, rideID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	if ride == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find ride with ID %s", id))
	}
	return c.JSON(http.StatusOK, ride)
}

func insertRide(db *sql.DB, ride Ride) (int64, error) {
	var placeholder = make([]string, len(tableColumns)-1)
	for i := range placeholder {
		placeholder[i] = "?"
	}

	stmt, err := db.Prepare("INSERT INTO rides (" + strings.Join(tableColumns[1:], ",") + ") VALUES (" + strings.Join(placeholder, ",") + ")")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		ride.StartLatitude,
		ride.StartLongitude,
		ride.EndLatitude,
		ride.EndLongitude,
		ride.RiderName,
		ride.DriverName,
		ride.DriverVehicle,
	)

	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func selectAllRides(db *sql.DB) ([]Ride, error) {
	rows, err := db.Query("SELECT " + strings.Join(tableColumns, ",") + " FROM rides")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rides := make([]Ride, 0)
	for rows.Next() {
		var ride Ride
		if err := rows.Scan(
			&ride.ID,
			&ride.StartLatitude,
			&ride.StartLongitude,
			&ride.EndLatitude,
			&ride.EndLongitude,
			&ride.RiderName,
			&ride.DriverName,
			&ride.DriverVehicle,
		); err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}

	return rides, nil
}

func selectRideByID(db *sql.DB, id int64) (*Ride, error) {
	var ride Ride
	err := db.QueryRow("SELECT "+strings.Join(tableColumns, ",")+" FROM rides WHERE id = ?", id).Scan(
		&ride.ID,
		&ride.StartLatitude,
		&ride.StartLongitude,
		&ride.EndLatitude,
		&ride.EndLongitude,
		&ride.RiderName,
		&ride.DriverName,
		&ride.DriverVehicle,
	)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return &ride, err
}
