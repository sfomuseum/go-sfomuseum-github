package messenger

import (
	"context"
	"fmt"
)

type MultiDeliveryAgent struct {
	DeliveryAgent
	agents []DeliveryAgent
}

func NewMultiDeliveryAgentWithURIs(ctx context.Context, agent_uris ...string) (DeliveryAgent, error) {

	agents := make([]DeliveryAgent, len(agent_uris))

	for idx, uri := range agent_uris {

		agent, err := NewDeliveryAgent(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create agent for %s, %w", uri, err)
		}

		agents[idx] = agent
	}

	return NewMultiDeliveryAgent(ctx, agents...)
}

func NewMultiDeliveryAgent(ctx context.Context, agents ...DeliveryAgent) (DeliveryAgent, error) {

	a := &MultiDeliveryAgent{
		agents: agents,
	}

	return a, nil
}

func (a *MultiDeliveryAgent) DeliverMessage(ctx context.Context, msg *Message) error {

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, agent := range a.agents {

		go func(agent DeliveryAgent) {

			err := agent.DeliverMessage(ctx, msg)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to deliver message for %T agent, %w", agent, err)
			}

			done_ch <- true

		}(agent)
	}

	remaining := len(a.agents)

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return err
		}
	}

	return nil
}
