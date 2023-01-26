package cronjob

import (
	"context"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gitlab.com/telegram_clone/chat_service/config"
	"gitlab.com/telegram_clone/chat_service/genproto/notification_service"
	grpcPkg "gitlab.com/telegram_clone/chat_service/pkg/grpc_client"
	"gitlab.com/telegram_clone/chat_service/storage"
	"gitlab.com/telegram_clone/chat_service/storage/repo"
)

type Cronjob struct {
	storage    storage.StorageI
	grpcClient grpcPkg.GrpcClientI
	cfg        *config.Config
	logger     *logrus.Logger
	cron       *cron.Cron
}

func NewCronjob(strg storage.StorageI, grpcConn grpcPkg.GrpcClientI, cfg *config.Config, logger *logrus.Logger) *Cronjob {
	c := cron.New()

	return &Cronjob{
		storage:    strg,
		grpcClient: grpcConn,
		cfg:        cfg,
		logger:     logger,
		cron:       c,
	}
}

func (cr *Cronjob) RegisterTasks() {
	cr.cron.AddFunc("@daily", cr.SendEmails)

	cr.cron.Start()
}

func (cr *Cronjob) SendEmails() {
	cr.logger.Info("******Sending daily emails*******")
	resp, err := cr.storage.User().GetAll(&repo.GetAllUsersParams{
		Page:  1,
		Limit: 1000,
	})
	if err != nil {
		cr.logger.Errorf("failed to get users: %v", err)
		return
	}

	for _, user := range resp.Users {
		_, err = cr.grpcClient.NotificationService().SendEmail(context.Background(), &notification_service.SendEmailRequest{
			To:      user.Email,
			Subject: "Daily news",
			Body: map[string]string{
				"title":       "Where does it come from?",
				"description": "Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of \"de Finibus Bonorum et Malorum\" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, \"Lorem ipsum dolor sit amet..\", comes from a line in section 1.10.32.",
			},
			Type: "news_email",
		})
		if err != nil {
			cr.logger.Errorf("failed to send email: %v", err)
			continue
		}
	}
}
