package gohassapi_test

import (
	ghs "github.com/korfuri/go-hass/pkg/gohassapi"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testingBearerToken = "1234abcd"
// func TestApiCheck(t *testing.T) {
// 	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 		assert.Equal(t, req.URL.String(), "/")
// 		rw.Write([]byte(`{"message":"All ok"}`))
// 	}))
// 	hc := ghs.NewClient(hs.URL, "token")
// 	r, err := hc.Check()
// 	assert.NoError(t, err)
// 	assert.Equal(t, r, "All ok")
// }

func timeMustParse(t *testing.T, timestr string) time.Time {
	ti, err := time.Parse(time.RFC3339Nano, timestr)
	assert.NoError(t, err, "Error while parsing a time.Time literal in test data")
	return ti
}

func TestGetCalls(t *testing.T) {
	for _, tc := range []struct {
		name string
		caller func(*testing.T, *ghs.HassClient) (any, error)
		httpResponse string
		httpStatusCode int
		expectedResponse any
	}{
		{
			name: "Check",
			caller: func(_ *testing.T, hc *ghs.HassClient) (any, error) {
					return hc.Check()
				},
			httpResponse: `{"message":"all OK"}`,
			httpStatusCode: 200,
			expectedResponse: "all OK",
		},
		{
			name: "States",
			caller: func(_ *testing.T, hc *ghs.HassClient) (any, error) {
				return hc.States()
			},
			httpResponse: `[
    {
        "attributes": {},
        "entity_id": "sun.sun",
        "last_changed": "2016-05-30T21:43:32.418320+00:00",
        "state": "below_horizon"
    },
    {
        "attributes": {},
        "entity_id": "process.Dropbox",
        "last_changed": "2016-05-30T21:43:32.418320+00:00",
        "state": "on"
    }
]`,
			httpStatusCode: 200,
			expectedResponse: []ghs.State{
				{
					Attributes: make(map[string]any),
					EntityId: "sun.sun",
					LastChanged: timeMustParse(t, "2016-05-30T21:43:32.418320+00:00"),
					State: "below_horizon",
				},
				{
					Attributes: make(map[string]any),
					EntityId: "process.Dropbox",
					LastChanged: timeMustParse(t, "2016-05-30T21:43:32.418320+00:00"),
					State: "on",
				},
			},
		},
		{
			name: "Services",
			caller: func(_ *testing.T, hc *ghs.HassClient) (any, error) {
				return hc.Services()
			},
			// TODO: Services should have a 'fields' field
			// but this is not well-documented in the API
			// reference, see:
			// https://github.com/home-assistant/developers.home-assistant/issues/2419
			httpResponse: `[
    {
        "domain": "test_domain",
        "services": {
           "test1": {
              "name": "test1",
              "description": "This is a test service"
            }
        }
    }
]`,
			httpStatusCode: 200,
			expectedResponse: []ghs.ServiceDomain{
				{
					Domain: "test_domain",
					Services: map[string]ghs.Service{
						"test1": {
							Name: "test1",
							Description: "This is a test service",
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var called bool
			hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				called = true
				assert.Equal(t, []string{fmt.Sprintf("Bearer %s", testingBearerToken)}, req.Header["Authorization"])
				assert.Equal(t, []string{"application/json"}, req.Header["Content-Type"])
				rw.Write([]byte(tc.httpResponse))
			}))
			hc := ghs.NewClient(hs.URL, testingBearerToken)
			r, err := tc.caller(t, hc)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResponse, r)
			assert.True(t, called, "The server method was not called at all")
		})
	}
}
