package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	v1 "gitlab.com/telegram_clone/api_gateway/api/v1"
	"gitlab.com/telegram_clone/api_gateway/config"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "gitlab.com/telegram_clone/api_gateway/api/docs" // for swagger

	grpcPkg "gitlab.com/telegram_clone/api_gateway/pkg/grpc_client"
)

type RouterOptions struct {
	Cfg        *config.Config
	GrpcClient grpcPkg.GrpcClientI
	Logger     *logrus.Logger
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is a blog service api.
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:        opt.Cfg,
		GrpcClient: opt.GrpcClient,
		Logger:     opt.Logger,
	})

	apiV1 := router.Group("/v1")

	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/verify", handlerV1.Verify)
	apiV1.POST("/auth/login", handlerV1.Login)
	apiV1.POST("/auth/forgot-password", handlerV1.ForgotPassword)
	apiV1.POST("/auth/verify-forgot-password", handlerV1.VerifyForgotPassword)
	apiV1.POST("/auth/update-password", handlerV1.AuthMiddleware("users", "update-password"), handlerV1.UpdatePassword)

	apiV1.POST("/users", handlerV1.AuthMiddleware("users", "create"), handlerV1.CreateUser)
	apiV1.GET("/users/:id", handlerV1.GetUser)
	apiV1.PUT("/users/:id", handlerV1.AuthMiddleware("users", "update"), handlerV1.UpdateUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	apiV1.DELETE("/users/:id", handlerV1.AuthMiddleware("users", "delete"), handlerV1.DeleteUser)
	apiV1.GET("/users/email/:email", handlerV1.GetUserByEmail)
	apiV1.GET("/users/me", handlerV1.AuthMiddleware("users", "get-profile"), handlerV1.GetUserByToken)

	apiV1.POST("/chats", handlerV1.AuthMiddleware("chats", "create"), handlerV1.CreateChat)
	apiV1.PUT("/chats/:id", handlerV1.AuthMiddleware("chats", "update"), handlerV1.UpdateChat)
	apiV1.DELETE("/chats/:id", handlerV1.AuthMiddleware("chats", "delete"), handlerV1.DeleteChat)
	apiV1.GET("/chats", handlerV1.AuthMiddleware("chats", "get"), handlerV1.GetAllChats)
	apiV1.GET("/chats/:id", handlerV1.AuthMiddleware("chats", "get-all"), handlerV1.GetChat)

	// group chat endpoints
	apiV1.POST("/chats/add-member", handlerV1.AuthMiddleware("chats", "add-member"), handlerV1.AddMember)
	apiV1.DELETE("/chats/remove-member", handlerV1.AuthMiddleware("chats", "remove-member"), handlerV1.RemoveMember)
	apiV1.DELETE("/chats/leave", handlerV1.AuthMiddleware("chats", "leave"), handlerV1.LeaveChat)
	apiV1.GET("/chats/members", handlerV1.AuthMiddleware("chats", "get-members"), handlerV1.GetChatMembers)

	apiV1.GET("/messages", handlerV1.GetAllMessages)

	apiV1.POST("/users/file-upload", handlerV1.AuthMiddleware("users", "users/file-upload"), handlerV1.UsersFileUpload)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
