package v1_test

import (
	"bytes"
	"context"
	"database/sql"
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
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCreateChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reqBody := models.ChatReq{
		Name:     faker.FirstName(),
		ChatType: "private",
		ImageUrl: faker.URL(),
	}

	Chat := mock_grpc.NewMockChatServiceClient(ctrl)
	Chat.EXPECT().Create(context.Background(), &pbc.Chat{
		Name:     reqBody.Name,
		ChatType: reqBody.ChatType,
		ImageUrl: reqBody.ImageUrl,
	}).Times(1).Return(&pbc.Chat{
		Id:   1,
		Name: reqBody.Name,
		UserInfo: &pbc.GetUserInfo{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Email:     faker.Email(),
			Username:  faker.Username(),
			ImageUrl:  faker.URL(),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		ChatType: reqBody.ChatType,
		ImageUrl: reqBody.ImageUrl,
	}, nil)

	payload, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	grpcConn.SetChatService(Chat)

	// TODO: after added permission it should be uncommented
	// accessToken := mockAuthMiddleware(t, ctrl)

	req, _ := http.NewRequest("POST", "/v1/chats", bytes.NewBuffer(payload))
	// req.Header.Add("Authorization", accessToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	body, _ := io.ReadAll(rec.Body)

	var response models.ChatRes
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
}

func TestGetChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chat := pbc.Chat{
		Id:     1,
		Name:   faker.FirstName(),
		UserId: 1,
		UserInfo: &pbc.GetUserInfo{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Email:     faker.Email(),
			Username:  faker.Username(),
			ImageUrl:  faker.URL(),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		ChatType: "private",
		ImageUrl: faker.URL(),
	}

	testCases := []struct {
		name          string
		Id            int64
		buildStubs    func(Chat *mock_grpc.MockChatServiceClient)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			Id:   chat.Id,
			buildStubs: func(Chat *mock_grpc.MockChatServiceClient) {
				Chat.EXPECT().Get(context.Background(), &pbc.ChatIdRequest{
					Id: chat.Id,
				}).Times(1).Return(&chat, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChat(t, recorder.Body, &chat)
			},
		},
		{
			name: "InternalError",
			Id:   chat.Id,
			buildStubs: func(Chat *mock_grpc.MockChatServiceClient) {
				Chat.EXPECT().Get(context.Background(), &pbc.ChatIdRequest{
					Id: chat.Id,
				}).Times(1).Return(&pbc.Chat{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	Chat := mock_grpc.NewMockChatServiceClient(ctrl)
	grpcConn.SetChatService(Chat)
	// TODO: after added permission it should be uncommented
	// accessToken := mockAuthMiddleware(t, ctrl)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(Chat)
			url := fmt.Sprintf("/v1/chats/%d", chat.Id)
			request, _ := http.NewRequest(http.MethodGet, url, nil)
			// request.Header.Add("Authorization", accessToken)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchChat(t *testing.T, body *bytes.Buffer, Chat *pbc.Chat) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var response models.ChatRes
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)
	assert.Equal(t, Chat.UserId, response.UserID)
}

func TestDeleteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ChatService := mock_grpc.NewMockChatServiceClient(ctrl)
	ChatService.EXPECT().Delete(context.Background(), &pbc.ChatIdRequest{
		Id: 1,
	}).Times(1).Return(&emptypb.Empty{}, nil)

	grpcConn.SetChatService(ChatService)

	// TODO: after added permission it should be uncommented
	// accessToken := mockAuthMiddleware(t, ctrl)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/chats/1", nil)
	// req.Header.Add("Authorization", accessToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, _ := io.ReadAll(rec.Body)

	var response models.ResponseOK
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Equal(t, response.Message, "success")
}
