package processor

import (
	"github.com/alexchebotarsky/heatpump-api/processor/event"
	"github.com/alexchebotarsky/heatpump-api/processor/handler"
	"github.com/alexchebotarsky/heatpump-api/processor/middleware"
)

func (p *Processor) setupEvents() {
	p.use(middleware.Metrics)

	p.handle(event.Event{
		Topic:   "heatpump/dht-sensor",
		Handler: handler.DHTSensor(p.Clients.Database),
	})
}
