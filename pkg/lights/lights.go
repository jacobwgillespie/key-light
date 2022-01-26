package lights

import (
	"context"
	"time"

	"github.com/endocrimes/keylight-go"
)

func DiscoverLights(ctx context.Context, timeoutDuration time.Duration) (<-chan *keylight.Device, error) {
	discovery, err := keylight.NewDiscovery()
	if err != nil {
		return nil, err
	}

	go func(d keylight.Discovery) {
		d.Run(ctx)
	}(discovery)

	return discovery.ResultsCh(), nil
}
