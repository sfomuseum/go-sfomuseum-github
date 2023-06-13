package messenger

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-slack"
)

type SlackDeliveryAgent struct {
	DeliveryAgent
	webhook *slack.Webhook
}

func init() {
	ctx := context.Background()
	RegisterDeliveryAgent(ctx, "slack", NewSlackDeliveryAgent)
}

func NewSlackDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	webhook_uri := q.Get("webhook")

	wh, err := slack.NewWebhook(ctx, webhook_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create webhook, %w", err)
	}

	a := &SlackDeliveryAgent{
		webhook: wh,
	}

	return a, nil
}

func (a *SlackDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {

	m := &slack.Message{
		Channel: msg.To,
		Text:    fmt.Sprintf("[%s] %s", msg.From, msg.Body),
	}

	a.webhook.Post(ctx, m)

	return nil
}
