package domain_test

import (
	"testing"

	domain "github.com/hawarir/backend-coding-test"
	"github.com/stretchr/testify/assert"
)

func TestRideValidation(t *testing.T) {
	testCases := []struct {
		testName    string
		ride        domain.Ride
		expectedErr string
	}{
		{
			testName: "When startLatitude and EndLatitude is incorrect",
			ride: domain.Ride{
				StartLatitude:  -91,
				StartLongitude: -180,
				EndLatitude:    91,
				EndLongitude:   180,
				RiderName:      "John Doe",
				DriverName:     "Driver",
				DriverVehicle:  "Car",
			},
			expectedErr: "-91.000000 is not a valid latitude value; 91.000000 is not a valid latitude value",
		},
		{
			testName: "When startLongitude and endLongitude is incorrect",
			ride: domain.Ride{
				StartLatitude:  -90,
				StartLongitude: -181,
				EndLatitude:    90,
				EndLongitude:   181,
				RiderName:      "John Doe",
				DriverName:     "Driver",
				DriverVehicle:  "Car",
			},
			expectedErr: "-181.000000 is not a valid longitude value; 181.000000 is not a valid longitude value",
		},
		{
			testName: "When riderName, driverName, and driverVehicle is empty",
			ride: domain.Ride{
				StartLatitude:  -90,
				StartLongitude: -180,
				EndLatitude:    90,
				EndLongitude:   180,
				RiderName:      "",
				DriverName:     "",
				DriverVehicle:  "",
			},
			expectedErr: "riderName can't be empty; driverName can't be empty; driverVehicle can't be empty",
		},
		{
			testName: "When values are correct",
			ride: domain.Ride{
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
			err := tc.ride.Validate()
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
