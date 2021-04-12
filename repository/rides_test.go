package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	domain "github.com/hawarir/backend-coding-test"
	"github.com/hawarir/backend-coding-test/repository"
)

type setupSQLMock func(sqlmock.Sqlmock)

func createRideRepo(fn setupSQLMock) (domain.RideRepository, *sql.DB) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if fn != nil {
		fn(mock)
	}
	return repository.NewRideRepository(db), db
}

func TestRideRepository_Insert(t *testing.T) {
	testCases := []struct {
		testName     string
		setupSQLMock setupSQLMock
		ride         domain.Ride
		lastInsertID int64
		expectedErr  string
	}{
		{
			testName: "When prepare returns error, return the error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO rides (startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle) VALUES (?,?,?,?,?,?,?)").
					WillReturnError(errors.New("Prepare error"))
			},
			ride:        domain.Ride{},
			expectedErr: "Prepare error",
		},
		{
			testName: "When exec returns error, return the error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO rides (startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle) VALUES (?,?,?,?,?,?,?)").
					WillReturnError(nil)

				mock.ExpectExec("INSERT INTO rides (startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle) VALUES (?,?,?,?,?,?,?)").
					WithArgs(
						float64(-90),
						float64(-180),
						float64(90),
						float64(180),
						"John Doe",
						"Driver",
						"Car",
					).WillReturnError(errors.New("Exec error"))
			},
			ride: domain.Ride{
				StartLatitude:  -90,
				StartLongitude: -180,
				EndLatitude:    90,
				EndLongitude:   180,
				RiderName:      "John Doe",
				DriverName:     "Driver",
				DriverVehicle:  "Car",
			},
			expectedErr: "Exec error",
		},
		{
			testName: "When successful, return the result",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO rides (startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle) VALUES (?,?,?,?,?,?,?)").
					WillReturnError(nil)

				mock.ExpectExec("INSERT INTO rides (startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle) VALUES (?,?,?,?,?,?,?)").
					WithArgs(
						float64(-90),
						float64(-180),
						float64(90),
						float64(180),
						"John Doe",
						"Driver",
						"Car",
					).WillReturnResult(sqlmock.NewResult(123, 1))
			},
			ride: domain.Ride{
				StartLatitude:  -90,
				StartLongitude: -180,
				EndLatitude:    90,
				EndLongitude:   180,
				RiderName:      "John Doe",
				DriverName:     "Driver",
				DriverVehicle:  "Car",
			},
			lastInsertID: 123,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rideRepo, db := createRideRepo(tc.setupSQLMock)
			defer db.Close()

			lastInsertID, err := rideRepo.Insert(tc.ride)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.lastInsertID, lastInsertID)
			}
		})
	}
}

func TestRideRepository_SelectAll(t *testing.T) {
	testCases := []struct {
		testName     string
		setupSQLMock setupSQLMock
		rides        []domain.Ride
		expectedErr  string
	}{
		{
			testName: "When query returns error, return the error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides").
					WillReturnError(errors.New("Query error"))
			},
			expectedErr: "Query error",
		},
		{
			testName: "When scan failed, return error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"id",
							"startLat",
							"startLong",
							"endLat",
							"endLong",
							"riderName",
							"driverName",
							"driverVehicle",
						}).
						AddRow(
							123,
							"not-a-number",
							"not-a-number",
							"not-a-number",
							"not-a-number",
							"John Doe",
							"Driver",
							"Car",
						))
			},
			expectedErr: "sql: Scan error on column index 1, name \"startLat\": converting driver.Value type string (\"not-a-number\") to a float64: invalid syntax",
		},
		{
			testName: "When return no rows, return empty slice",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"id",
							"startLat",
							"startLong",
							"endLat",
							"endLong",
							"riderName",
							"driverName",
							"driverVehicle",
						}))
			},
			rides: []domain.Ride{},
		},
		{
			testName: "When successful, return rides",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"id",
							"startLat",
							"startLong",
							"endLat",
							"endLong",
							"riderName",
							"driverName",
							"driverVehicle",
						}).
						AddRow(
							123,
							-90,
							-180,
							90,
							180,
							"John Doe",
							"Driver",
							"Car",
						))
			},
			rides: []domain.Ride{
				{
					ID:             123,
					StartLatitude:  -90,
					StartLongitude: -180,
					EndLatitude:    90,
					EndLongitude:   180,
					RiderName:      "John Doe",
					DriverName:     "Driver",
					DriverVehicle:  "Car",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rideRepo, db := createRideRepo(tc.setupSQLMock)
			defer db.Close()

			rides, err := rideRepo.SelectAll()
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.rides, rides)
			}
		})
	}
}

func TestRideRepository_SelectByID(t *testing.T) {
	testCases := []struct {
		testName     string
		setupSQLMock setupSQLMock
		rideID       int64
		ride         *domain.Ride
		expectedErr  string
	}{
		{
			testName: "When query returns error, return the error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides WHERE id = ?").
					WithArgs(int64(123)).
					WillReturnError(errors.New("Query error"))
			},
			rideID:      123,
			expectedErr: "Query error",
		},
		{
			testName: "When query returns errNoRows, return nil",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides WHERE id = ?").
					WithArgs(int64(123)).
					WillReturnError(sql.ErrNoRows)
			},
			rideID: 123,
			ride:   nil,
		},
		{
			testName: "When scan failed, return error",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides WHERE id = ?").
					WithArgs(int64(123)).
					WillReturnRows(sqlmock.
						NewRows([]string{
							"id",
							"startLat",
							"startLong",
							"endLat",
							"endLong",
							"riderName",
							"driverName",
							"driverVehicle",
						}).
						AddRow(
							123,
							"not-a-number",
							"not-a-number",
							"not-a-number",
							"not-a-number",
							"John Doe",
							"Driver",
							"Car",
						))
			},
			rideID:      123,
			expectedErr: "sql: Scan error on column index 1, name \"startLat\": converting driver.Value type string (\"not-a-number\") to a float64: invalid syntax",
		},
		{
			testName: "When successful, return ride",
			setupSQLMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, startLat, startLong, endLat, endLong, riderName, driverName, driverVehicle FROM rides WHERE id = ?").
					WithArgs(int64(123)).
					WillReturnRows(sqlmock.
						NewRows([]string{
							"id",
							"startLat",
							"startLong",
							"endLat",
							"endLong",
							"riderName",
							"driverName",
							"driverVehicle",
						}).
						AddRow(
							123,
							-90,
							-180,
							90,
							180,
							"John Doe",
							"Driver",
							"Car",
						))
			},
			rideID: 123,
			ride: &domain.Ride{
				ID:             123,
				StartLatitude:  -90,
				StartLongitude: -180,
				EndLatitude:    90,
				EndLongitude:   180,
				RiderName:      "John Doe",
				DriverName:     "Driver",
				DriverVehicle:  "Car",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			rideRepo, db := createRideRepo(tc.setupSQLMock)
			defer db.Close()

			ride, err := rideRepo.SelectByID(tc.rideID)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ride, ride)
			}
		})
	}
}
