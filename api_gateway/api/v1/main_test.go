package v1_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/telegram_clone/api_gateway/api"
	"gitlab.com/telegram_clone/api_gateway/config"
	grpcPkg "gitlab.com/telegram_clone/api_gateway/pkg/grpc_client"
	"gitlab.com/telegram_clone/api_gateway/pkg/logger"
)

var (
	router   *gin.Engine
	grpcConn grpcPkg.GrpcClientI
)

func TestMain(m *testing.M) {
	var err error
	cfg := config.Load("./../..")

	lgr := logger.New()

	grpcConn, err = grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}

	ginEngine := api.New(&api.RouterOptions{
		Cfg:        &cfg,
		GrpcClient: grpcConn,
		Logger:     lgr,
	})

	router = ginEngine
	os.Exit(m.Run())
}
