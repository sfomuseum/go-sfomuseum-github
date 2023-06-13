package messenger

import (
	"context"
)

type NullDeliveryAgent struct {
	DeliveryAgent
}

func init() {
	ctx := context.Background()
	RegisterDeliveryAgent(ctx, "null", NewNullDeliveryAgent)
}

func NewNullDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	a := &NullDeliveryAgent{}
	return a, nil
}

func (a *NullDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {
	return nil
}
