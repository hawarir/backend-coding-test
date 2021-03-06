package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	domain "github.com/hawarir/backend-coding-test"
)

type (
	rideCntrl struct {
		rideRepo domain.RideRepository
	}

	ridesEnvelope struct {
		Rides  []domain.Ride `json:"rides"`
		Cursor string        `json:"cursor"`
	}
)

func SetupRideController(e *echo.Echo, rideRepo domain.RideRepository) {
	cntrl := &rideCntrl{rideRepo: rideRepo}

	e.GET("/health", healthCheck)

	e.POST("/rides", cntrl.addRide)
	e.GET("/rides", cntrl.getAllRides)
	e.GET("/rides/:id", cntrl.getRide)
}

func healthCheck(c echo.Context) error {
	return c.String(200, "Healthy")
}

func (cntrl rideCntrl) addRide(c echo.Context) error {
	var ride domain.Ride
	if err := c.Bind(&ride); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Malformed request body: %s", err))
	}
	if err := ride.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid request body: %s", err))
	}
	lastInsertID, err := cntrl.rideRepo.Insert(ride)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	ride.ID = lastInsertID
	return c.JSON(http.StatusCreated, ride)
}

func (cntrl rideCntrl) getAllRides(c echo.Context) error {
	var page domain.Pagination
	if err := c.Bind(&page); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Bad request: %s", err))
	}
	rides, cursor, err := cntrl.rideRepo.SelectAll(page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	return c.JSON(http.StatusOK, ridesEnvelope{Rides: rides, Cursor: cursor})
}

func (cntrl rideCntrl) getRide(c echo.Context) error {
	id := c.Param("id")
	rideID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid ID: %s", err))
	}
	ride, err := cntrl.rideRepo.SelectByID(rideID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", err))
	}
	if ride == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Can't find ride with ID %s", id))
	}
	return c.JSON(http.StatusOK, ride)
}
