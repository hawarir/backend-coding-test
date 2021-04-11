package domain

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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

	RideRepository interface {
		InitTable() error

		Insert(Ride) (int64, error)
		SelectAll() ([]Ride, error)
		SelectByID(int64) (*Ride, error)
	}
)

func (r Ride) Validate() error {
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
