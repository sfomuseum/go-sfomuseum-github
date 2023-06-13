package messenger

import (
	"context"
	"fmt"
)

type StdoutDeliveryAgent struct {
	DeliveryAgent
}

func init() {
	ctx := context.Background()
	RegisterDeliveryAgent(ctx, "stdout", NewStdoutDeliveryAgent)
}

func NewStdoutDeliveryAgent(ctx context.Context, uri string) (DeliveryAgent, error) {

	a := &StdoutDeliveryAgent{}
	return a, nil
}

func (a *StdoutDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {
	fmt.Println(msg.Body)
	return nil
}
