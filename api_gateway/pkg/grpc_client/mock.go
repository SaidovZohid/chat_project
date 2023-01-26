package grpc_client

import pbu "gitlab.com/telegram_clone/api_gateway/genproto/chat_service"

func (g *GrpcClient) SetUserService(u pbu.UserServiceClient) {
	g.connections["user_service"] = u
}

func (g *GrpcClient) SetAuthService(a pbu.AuthServiceClient) {
	g.connections["auth_service"] = a
}

func (g *GrpcClient) SetChatService(p pbu.ChatServiceClient) {
	g.connections["chat_service"] = p
}

func (g *GrpcClient) SetMessageService(p pbu.MessageServiceClient) {
	g.connections["message_service"] = p
}
