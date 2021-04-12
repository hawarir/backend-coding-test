package repository

import (
	"database/sql"
	"fmt"
	"strings"

	domain "github.com/hawarir/backend-coding-test"
)

type rideRepository struct {
	db              *sql.DB
	tableColumns    []string
	tableDefinition []string
}

func NewRideRepository(db *sql.DB) domain.RideRepository {
	tableSchema := [][2]string{
		{"id", "INTEGER PRIMARY KEY AUTOINCREMENT"},
		{"startLat", "REAL NOT NULL"},
		{"startLong", "REAL NOT NULL"},
		{"endLat", "REAL NOT NULL"},
		{"endLong", "REAL NOT NULL"},
		{"riderName", "TEXT NOT NULL"},
		{"driverName", "TEXT NOT NULL"},
		{"driverVehicle", "TEXT NOT NULL"},
	}

	tableColumns := make([]string, len(tableSchema))
	tableDefinition := make([]string, len(tableSchema))

	for i, tuple := range tableSchema {
		tableColumns[i] = tuple[0]
		tableDefinition[i] = fmt.Sprintf("%s %s", tuple[0], tuple[1])
	}
	return rideRepository{db: db, tableColumns: tableColumns, tableDefinition: tableDefinition}
}

// NOTE: This shouldn't be needed in production environment
func (r rideRepository) InitTable() error {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS rides (" + strings.Join(r.tableDefinition, ",") + ")")
	return err
}

func (r rideRepository) Insert(ride domain.Ride) (int64, error) {
	var placeholder = make([]string, len(r.tableColumns)-1)
	for i := range placeholder {
		placeholder[i] = "?"
	}

	stmt, err := r.db.Prepare("INSERT INTO rides (" + strings.Join(r.tableColumns[1:], ",") + ") VALUES (" + strings.Join(placeholder, ",") + ")")
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
	rows, err := r.db.Query("SELECT " + strings.Join(r.tableColumns, ",") + " FROM rides")
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
	err := r.db.QueryRow("SELECT "+strings.Join(r.tableColumns, ",")+" FROM rides WHERE id = ?", id).Scan(
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
