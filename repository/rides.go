package repository

import (
	"database/sql"
	"fmt"
	"strings"

	domain "github.com/hawarir/backend-coding-test"
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

type rideRepository struct {
	db *sql.DB
}

func NewRideRepository(db *sql.DB) domain.RideRepository {
	return rideRepository{db: db}
}

// NOTE: This shouldn't be needed in production environment
func (r rideRepository) InitTable() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS rides (" + strings.Join(tableDefinition, ",") + ")")
	return err
}

func (r rideRepository) Insert(ride domain.Ride) (int64, error) {
	var placeholder = make([]string, len(tableColumns)-1)
	for i := range placeholder {
		placeholder[i] = "?"
	}

	stmt, err := r.db.Prepare("INSERT INTO rides (" + strings.Join(tableColumns[1:], ",") + ") VALUES (" + strings.Join(placeholder, ",") + ")")
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

func (r rideRepository) SelectAll() ([]domain.Ride, error) {
	rows, err := r.db.Query("SELECT " + strings.Join(tableColumns, ",") + " FROM rides")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rides := make([]domain.Ride, 0)
	for rows.Next() {
		var ride domain.Ride
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

func (r rideRepository) SelectByID(id int64) (*domain.Ride, error) {
	var ride domain.Ride
	err := r.db.QueryRow("SELECT "+strings.Join(tableColumns, ",")+" FROM rides WHERE id = ?", id).Scan(
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
