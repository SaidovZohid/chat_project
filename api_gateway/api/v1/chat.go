package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/telegram_clone/api_gateway/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "gitlab.com/telegram_clone/api_gateway/genproto/chat_service"
)

// @Security ApiKeyAuth
// @Router /chats [post]
// @Summary Create Chat
// @Description Create Chat
// @Tags chat
// @Accept json
// @Produce json
// @Param user body models.ChatReq true "User"
// @Success 201 {object} models.Chat
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) CreateChat(c *gin.Context) {
	var (
		req models.ChatReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	chat, err := h.grpcClient.ChatService().Create(context.Background(), &pbc.CreateChatReq{
		Name:     req.Name,
		UserId:   payload.UserID,
		ChatType: req.ChatType,
		ImageUrl: req.ImageUrl,
		Members:  req.Members,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to create chat")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseChat(chat))
}

func parseChat(chat *pbc.Chat) models.Chat {
	return models.Chat{
		ID:     chat.Id,
		Name:   chat.Name,
		UserID: chat.UserId,
		UserInfo: models.GetUserInfo{
			FirstName: chat.UserInfo.FirstName,
			LastName:  chat.UserInfo.LastName,
			Email:     chat.UserInfo.Email,
			Username:  chat.UserInfo.Username,
			ImageUrl:  chat.UserInfo.ImageUrl,
			CreatedAt: chat.UserInfo.CreatedAt,
		},
		ChatType: chat.ChatType,
		ImageUrl: chat.ImageUrl,
	}
}

// @Security ApiKeyAuth
// @Router /chats/{id} [get]
// @Summary Get chat and users info
// @Description Get chat and users info
// @Tags chat
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} models.Chat
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
func (h *handlerV1) GetChat(c *gin.Context) {
	h.logger.Info("get chat info")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	chat, err := h.grpcClient.ChatService().Get(context.Background(), &pbc.IdRequest{
		Id: id,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to get chat info")

		if s, _ := status.FromError(err); s.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseChat(chat))
}

// @Security ApiKeyAuth
// @Router /chats/{id} [delete]
// @Summary Delete chat
// @Description Delete chat
// @Tags chat
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
func (h *handlerV1) DeleteChat(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.grpcClient.ChatService().Delete(context.Background(), &pbc.ChatIdRequest{
		Id: id,
	})
	if err != nil {
		if s, _ := status.FromError(err); s.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /chats [get]
// @Summary Get all chats
// @Description Get all chats
// @Tags chat
// @Accept json
// @Produce json
// @Param filter query models.GetAllChatsParams false "Filter"
// @Success 200 {object} models.GetAllChatsRes
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) GetAllChats(c *gin.Context) {
	req, err := validateGetAllChatsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.ChatService().GetAll(context.Background(), &pbc.GetAllChatsParams{
		Page:   req.Page,
		Limit:  req.Limit,
		UserId: req.UserID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to get all chats")
		if s, _ := status.FromError(err); s.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := models.GetAllChatsRes{
		Chats: make([]*models.Chat, 0),
		Count: result.Count,
	}
	for _, v := range result.Chats {
		res := parseChat(v)
		response.Chats = append(response.Chats, &res)

	}
	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /chats [put]
// @Summary Update Chat
// @Description Update Chat
// @Tags chat
// @Accept json
// @Produce json
// @Param user body models.ChatReq true "User"
// @Success 201 {object} models.Chat
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) UpdateChat(c *gin.Context) {
	var (
		req models.Chat
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	chat, err := h.grpcClient.ChatService().Update(context.Background(), &pbc.Chat{
		Name:     req.Name,
		UserId:   req.UserID,
		ChatType: req.ChatType,
		ImageUrl: req.ImageUrl,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update chat")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseChat(chat))
}

// @Security ApiKeyAuth
// @Router /chats/add-member [post]
// @Summary Add member to group chat
// @Description Add member to group chat
// @Tags chat
// @Accept json
// @Produce json
// @Param data body models.AddRemoveMemberReq true "data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) AddMember(c *gin.Context) {
	var (
		req models.AddRemoveMemberReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := h.grpcClient.ChatService().AddMember(context.Background(), &pbc.AddMemberRequest{
		UserId: req.UserID,
		ChatId: req.ChatID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update chat")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /chats/remove-member [delete]
// @Summary Remove member from group chat
// @Description Remove member from group chat
// @Tags chat
// @Accept json
// @Produce json
// @Param data body models.AddRemoveMemberReq true "data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) RemoveMember(c *gin.Context) {
	var (
		req models.AddRemoveMemberReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := h.grpcClient.ChatService().RemoveMember(context.Background(), &pbc.RemoveMemberRequest{
		UserId: req.UserID,
		ChatId: req.ChatID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update chat")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /chats/leave [delete]
// @Summary Leave from group chat
// @Description Leave from group chat
// @Tags chat
// @Accept json
// @Produce json
// @Param data body models.LeaveGroupReq true "data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
func (h *handlerV1) LeaveChat(c *gin.Context) {
	var (
		req models.LeaveGroupReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = h.grpcClient.ChatService().RemoveMember(context.Background(), &pbc.RemoveMemberRequest{
		UserId: payload.UserID,
		ChatId: req.ChatID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to update chat")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "success",
	})
}

// @Security ApiKeyAuth
// @Router /chats/members [get]
// @Summary Get all chat members
// @Description Get all chat members
// @Tags chat
// @Accept json
// @Produce json
// @Param filter query models.GetChatMembersParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetChatMembers(c *gin.Context) {
	req, err := validateGetChatMembersParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.ChatService().GetChatMembers(context.Background(), &pbc.GetChatMembersParams{
		Page:   req.Page,
		Limit:  req.Limit,
		ChatId: req.ChatID,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to get all users")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

func validateGetChatMembersParams(c *gin.Context) (*models.GetChatMembersParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	chatID, err := strconv.Atoi(c.Query("chat_id"))
	if err != nil {
		return nil, err
	}

	return &models.GetChatMembersParams{
		Limit:  int64(limit),
		Page:   int64(page),
		ChatID: int64(chatID),
	}, nil
}
