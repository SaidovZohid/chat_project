package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func loginUser(t *testing.T) *models.AuthResponse {
	resp := httptest.NewRecorder()

	payload, err := json.Marshal(models.LoginRequest{
		Email:    "testuser@gmail.com",
		Password: "asdf1234",
	})
	assert.NoError(t, err)
	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(payload))
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	body, _ := io.ReadAll(resp.Body)

	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	return &response
}

func loginSuperadmin(t *testing.T) *models.AuthResponse {
	resp := httptest.NewRecorder()

	payload, err := json.Marshal(models.LoginRequest{
		Email:    "t.mannonov@gmail.com",
		Password: "asdf1234",
	})
	assert.NoError(t, err)
	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(payload))
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	body, _ := io.ReadAll(resp.Body)

	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	return &response
}

func TestLogin(t *testing.T) {
	loginUser(t)
	loginSuperadmin(t)
}

func TestLoginMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.LoginRequest{
		Email:    "t.mannonov@gmail.com",
		Password: "asdf1234",
	}

	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().Login(context.Background(), &pbc.LoginRequest{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}).Times(1).Return(&pbc.AuthResponse{
		Id:          1,
		FirstName:   "Temur",
		LastName:    "Mannonov",
		Email:       "t.mannonov@gmail.com",
		Username:    "temur",
		Type:        "superadmin",
		CreatedAt:   time.Now().Format(time.RFC3339),
		AccessToken: faker.Sentence(),
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetAuthService(authService)

	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)

	body, _ := io.ReadAll(rec.Body)

	var response models.AuthResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
}

func mockAuthMiddleware(t *testing.T, ctrl *gomock.Controller) string {
	accessToken := faker.UUIDHyphenated()

	// mocking auth
	authService := mock_grpc.NewMockAuthServiceClient(ctrl)
	authService.EXPECT().VerifyToken(context.Background(), &pbc.VerifyTokenRequest{
		AccessToken: accessToken,
		Resource:    "users",
		Action:      "create",
	}).Times(1).Return(&pbc.AuthPayload{
		Id:            faker.UUIDHyphenated(),
		UserId:        1,
		Email:         faker.Email(),
		UserType:      "superadmin",
		HasPermission: true,
		IssuedAt:      time.Now().Format(time.RFC3339),
		ExpiredAt:     time.Now().Add(time.Hour).Format(time.RFC3339),
	}, nil)

	grpcConn.SetAuthService(authService)

	return accessToken
}
