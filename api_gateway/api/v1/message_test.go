package v1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllMessages(t *testing.T) {
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/messages", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	fmt.Println(resp.Body.String())
}

func TestGetAllMessagesCases(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "success case",
			query: "?limit=10&page=1",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name:  "incorrect limit param",
			query: "?limit=ads&page=1",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/messages%s", tc.query)

			req, _ := http.NewRequest("GET", url, nil)
			router.ServeHTTP(resp, req)

			tc.checkResponse(t, resp)
		})
	}
}
