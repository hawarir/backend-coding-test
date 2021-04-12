package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	domain "github.com/hawarir/backend-coding-test"
	"github.com/hawarir/backend-coding-test/repository/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type setupMockRepo func(mockRepo *mock.MockRideRepository)

func NewRideController(t *testing.T, fn setupMockRepo) (rideCntrl, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)

	rideRepo := mock.NewMockRideRepository(mockCtrl)
	if fn != nil {
		fn(rideRepo)
	}

	return rideCntrl{rideRepo: rideRepo}, mockCtrl
}

func TestRideController_addRide(t *testing.T) {
	testCases := []struct {
		testName      string
		requestBody   string
		setupMockRepo setupMockRepo
		statusCode    int
		responseBody  string
		expectedErr   string
	}{
		{
			testName:    "When request body is malformed, return status code 400 with error message",
			requestBody: "invalid-json",
			statusCode:  http.StatusBadRequest,
			expectedErr: "code=400, message=Malformed request body: code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value",
		},
		{
			testName:    "When request is invalid, return status code 422 with error message",
			requestBody: `{"startLatitude": -100, "startLongitude": -200, "endLatitude": 100, "endLongitude": 200, "riderName": "", "driverName": "", "driverVehicle": ""}`,
			statusCode:  http.StatusUnprocessableEntity,
			expectedErr: "code=422, message=Invalid request body: -100.000000 is not a valid latitude value; 100.000000 is not a valid latitude value; -200.000000 is not a valid longitude value; 200.000000 is not a valid longitude value; riderName can't be empty; driverName can't be empty; driverVehicle can't be empty",
		},
		{
			testName:    "When repository returns error, return status code 500 with error message",
			requestBody: `{"startLatitude": 90, "startLongitude": 180, "endLatitude": 90, "endLongitude": 180, "riderName": "John Doe", "driverName": "Driver", "driverVehicle": "Car"}`,
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().
					Insert(domain.Ride{
						StartLatitude:  90,
						StartLongitude: 180,
						EndLatitude:    90,
						EndLongitude:   180,
						RiderName:      "John Doe",
						DriverName:     "Driver",
						DriverVehicle:  "Car",
					}).
					Return(int64(-1), errors.New("Insert error"))
			},
			statusCode:  http.StatusInternalServerError,
			expectedErr: "code=500, message=Internal server error: Insert error",
		},
		{
			testName:    "When successful, return status code 201 with response body",
			requestBody: `{"startLatitude": 90, "startLongitude": 180, "endLatitude": 90, "endLongitude": 180, "riderName": "John Doe", "driverName": "Driver", "driverVehicle": "Car"}`,
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().
					Insert(domain.Ride{
						StartLatitude:  90,
						StartLongitude: 180,
						EndLatitude:    90,
						EndLongitude:   180,
						RiderName:      "John Doe",
						DriverName:     "Driver",
						DriverVehicle:  "Car",
					}).
					Return(int64(1), nil)
			},
			statusCode:   http.StatusCreated,
			responseBody: "{\"id\":1,\"startLatitude\":90,\"startLongitude\":180,\"endLatitude\":90,\"endLongitude\":180,\"riderName\":\"John Doe\",\"driverName\":\"Driver\",\"driverVehicle\":\"Car\"}\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/rides", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			cntrl, mock := NewRideController(t, tc.setupMockRepo)
			defer mock.Finish()

			err := cntrl.addRide(c)
			if tc.expectedErr != "" {
				httpErr, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tc.statusCode, httpErr.Code)
					assert.Equal(t, tc.expectedErr, err.Error())
				}
			} else {
				assert.Equal(t, tc.statusCode, rec.Code)
				assert.Equal(t, tc.responseBody, rec.Body.String())
			}
		})
	}
}

func TestRideController_getAllRides(t *testing.T) {
	testCases := []struct {
		testName      string
		setupMockRepo setupMockRepo
		statusCode    int
		responseBody  string
		expectedErr   string
	}{
		{
			testName: "When repository returns error, return status code 500 with error message",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().SelectAll().
					Return(nil, errors.New("Select All error"))
			},
			statusCode:  http.StatusInternalServerError,
			expectedErr: "code=500, message=Internal server error: Select All error",
		},
		{
			testName: "When repository returns empty result, return status code 200 with empty array in response body",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().SelectAll().
					Return([]domain.Ride{}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: "[]\n",
		},
		{
			testName: "When repository returns results, return status code 200 with the results as array",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().SelectAll().
					Return([]domain.Ride{
						{
							ID:             1,
							StartLatitude:  90,
							StartLongitude: 180,
							EndLatitude:    90,
							EndLongitude:   180,
							RiderName:      "John Doe",
							DriverName:     "Driver",
							DriverVehicle:  "Car",
						},
					}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: "[{\"id\":1,\"startLatitude\":90,\"startLongitude\":180,\"endLatitude\":90,\"endLongitude\":180,\"riderName\":\"John Doe\",\"driverName\":\"Driver\",\"driverVehicle\":\"Car\"}]\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/rides", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			cntrl, mock := NewRideController(t, tc.setupMockRepo)
			defer mock.Finish()

			err := cntrl.getAllRides(c)
			if tc.expectedErr != "" {
				httpErr, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tc.statusCode, httpErr.Code)
					assert.Equal(t, tc.expectedErr, err.Error())
				}
			} else {
				assert.Equal(t, tc.statusCode, rec.Code)
				assert.Equal(t, tc.responseBody, rec.Body.String())
			}
		})
	}
}

func TestRideController_getRide(t *testing.T) {
	testCases := []struct {
		testName      string
		paramID       string
		setupMockRepo setupMockRepo
		statusCode    int
		responseBody  string
		expectedErr   string
	}{
		{
			testName:    "When ID is not an integer, return status code 422 with error message",
			paramID:     "not-a-string",
			statusCode:  http.StatusUnprocessableEntity,
			expectedErr: "code=422, message=Invalid ID: strconv.ParseInt: parsing \"not-a-string\": invalid syntax",
		},
		{
			testName: "When repository returns error, return status code 500 with error message",
			paramID:  "1",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().
					SelectByID(int64(1)).
					Return(nil, errors.New("Select By ID error"))
			},
			statusCode:  http.StatusInternalServerError,
			expectedErr: "code=500, message=Internal server error: Select By ID error",
		},
		{
			testName: "When repository returns no result, return status code 404 with error message",
			paramID:  "1",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().
					SelectByID(int64(1)).
					Return(nil, nil)
			},
			statusCode:  http.StatusNotFound,
			expectedErr: "code=404, message=Can't find ride with ID 1",
		},
		{
			testName: "When successful, return status code 200 with result",
			paramID:  "1",
			setupMockRepo: func(mockRepo *mock.MockRideRepository) {
				mockRepo.EXPECT().
					SelectByID(int64(1)).
					Return(&domain.Ride{
						ID:             1,
						StartLatitude:  90,
						StartLongitude: 180,
						EndLatitude:    90,
						EndLongitude:   180,
						RiderName:      "John Doe",
						DriverName:     "Driver",
						DriverVehicle:  "Car",
					}, nil)
			},
			statusCode:   http.StatusOK,
			responseBody: "{\"id\":1,\"startLatitude\":90,\"startLongitude\":180,\"endLatitude\":90,\"endLongitude\":180,\"riderName\":\"John Doe\",\"driverName\":\"Driver\",\"driverVehicle\":\"Car\"}\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/rides/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.paramID)

			cntrl, mock := NewRideController(t, tc.setupMockRepo)
			defer mock.Finish()

			err := cntrl.getRide(c)
			if tc.expectedErr != "" {
				httpErr, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tc.statusCode, httpErr.Code)
					assert.Equal(t, tc.expectedErr, err.Error())
				}
			} else {
				assert.Equal(t, tc.statusCode, rec.Code)
				assert.Equal(t, tc.responseBody, rec.Body.String())
			}
		})
	}
}
