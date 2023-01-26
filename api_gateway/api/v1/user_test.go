package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/telegram_clone/api_gateway/api/models"
	pbc "gitlab.com/telegram_clone/api_gateway/genproto/chat_service"
	"gitlab.com/telegram_clone/api_gateway/pkg/grpc_client/mock_grpc"
)

func TestGetAllUsers(t *testing.T) {
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users", nil)
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	fmt.Println(resp.Body.String())
}

func TestGetAllUsersCases(t *testing.T) {
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
			url := fmt.Sprintf("/v1/users%s", tc.query)

			req, _ := http.NewRequest("GET", url, nil)
			router.ServeHTTP(resp, req)

			tc.checkResponse(t, resp)
		})
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.CreateUserRequest{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  "@Qwerty123",
		Type:      "superadmin",
	}

	userService := mock_grpc.NewMockUserServiceClient(ctrl)
	userService.EXPECT().Create(context.Background(), &pbc.User{
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Email:     reqBody.Email,
		Password:  reqBody.Password,
		Type:      reqBody.Type,
	}).Times(1).Return(&pbc.User{
		Id:        1,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Email:     reqBody.Email,
		Type:      reqBody.Type,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetUserService(userService)

	accessToken := mockAuthMiddleware(t, ctrl)

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	req.Header.Add("Authorization", accessToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	body, _ := io.ReadAll(rec.Body)

	var response models.User
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.FirstName, response.FirstName)
	assert.Equal(t, reqBody.Email, response.Email)
}
