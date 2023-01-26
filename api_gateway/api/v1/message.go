package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/telegram_clone/api_gateway/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "gitlab.com/telegram_clone/api_gateway/genproto/chat_service"
)

func parseMessage(message *pbc.ChatMessage) models.Message {
	return models.Message{
		ID:      message.Id,
		Message: message.Message,
		UserID: message.UserId,
		UserInfo: models.GetUserInfo{
			FirstName: message.UserInfo.FirstName,
			LastName:  message.UserInfo.LastName,
			Email:     message.UserInfo.Email,
			Username:  message.UserInfo.Username,
			ImageUrl:  message.UserInfo.ImageUrl,
			CreatedAt: message.UserInfo.CreatedAt,
		},
		ChatID: message.ChatId,
		CreatedAt: message.CreatedAt,
	}
}



// @Router /messages [get]
// @Summary Get all messages
// @Description Get all messages
// @Tags message
// @Accept json
// @Produce json
// @Param filter query models.GetAllMessagesParams false "Filter"
// @Success 200 {object} models.GetAllMessagesRes
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) GetAllMessages(c *gin.Context) {
	req, err := validateGetAllMessagesParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.MessageService().GetAll(context.Background(), &pbc.GetAllMessagesParams{
		Page:  req.Page,
		Limit: req.Limit,
		ChatId: req.ChatID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to get all messages")
		if s, _ := status.FromError(err); s.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := models.GetAllMessagesRes{
		Messages: make([]*models.Message, 0),
		Count: result.Count,
	}
	for _, v := range result.Messages {
		res := parseMessage(v)
		response.Messages = append(response.Messages, &res)

	}
	c.JSON(http.StatusOK, response)
}

