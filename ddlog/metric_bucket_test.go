package ddlog

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func TestMetricStore(t *testing.T) {
	config := NewConfig()
	config.AlertDuration = 2 * time.Minute
	config.AlertThreshold = 100
	config.FastTickDuration = 10 * time.Second
	ms := NewMetricStore(config)
	ms.clock = clock.NewMock()
	g := NewGenerator(nil)

	msgChan := make(chan *Message, 10)
	for i := 0; i < 10; i++ {
		msgChan <- g.RandMsg()
	}

	ms.Start(msgChan, func(e *MetricEvent) {

	})

	// alpha := alphaVal(1*time.Second, 2*time.Minute))
}
