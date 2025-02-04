package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
)

type IRSignal struct {
	Signal string `json:"signal"`
}

func (p *PubSub) TransmitIRSignal(ctx context.Context, binaryString string) error {
	signal := IRSignal{
		Signal: binaryString,
	}
	payload, err := json.Marshal(&signal)
	if err != nil {
		return fmt.Errorf("error marshalling ir signal: %v", err)
	}

	err = p.Publish(ctx, "heatpump/ir-transmitter", payload)
	if err != nil {
		return fmt.Errorf("error publishing heatpump binary state: %v", err)
	}

	return nil
}
