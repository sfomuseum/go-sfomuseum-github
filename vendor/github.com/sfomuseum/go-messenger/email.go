package messenger

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/gomail-sender"
	"github.com/aaronland/gomail/v2"
)

type EmailDeliveryAgent struct {
	DeliveryAgent
	sender gomail.Sender
}

func init() {

	ctx := context.Background()

	for _, uri := range sender.Schemes() {
		s := strings.Replace(uri, "://", "", 1)
		s = fmt.Sprintf("email-%s", s)
		RegisterDeliveryAgent(ctx, s, NewEmailDeliveryAgent)
	}
}

func NewEmailDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	s, err := sender.NewSender(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new email sender, %w", err)
	}

	agent := &EmailDeliveryAgent{
		sender: s,
	}

	return agent, nil
}

func (agent *EmailDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {

	gomail_msg := gomail.NewMessage()

	gomail_msg.SetHeader("Subject", msg.Subject)
	gomail_msg.SetHeader("From", msg.From)
	gomail_msg.SetHeader("To", msg.To)

	gomail_msg.SetBody("text/html", msg.Body)

	err := gomail.Send(agent.sender, gomail_msg)

	if err != nil {
		return fmt.Errorf("Failed to deliver message, %w", err)
	}

	return nil
}
