package mailer

import (
	"github.com/hillview.tv/coreAPI/env"
	"github.com/sendgrid/sendgrid-go"
)

func GetSendgridClient() *sendgrid.Client {
	client := sendgrid.NewSendClient(env.SendgridAPIKey)
	return client
}
