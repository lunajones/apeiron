package world

import (
	"time"

	"github.com/lunajones/apeiron/service/zone"
)

func TickAll() {
	const targetFPS = 60
	const targetDelta = 1.0 / float64(targetFPS)

	var lastTick = time.Now()

	for {
		now := time.Now()
		elapsed := now.Sub(lastTick).Seconds()
		if elapsed > 0.1 {
			elapsed = 0.1
		}
		lastTick = now

		for _, z := range zone.Zones {
			z.Tick(elapsed)
		}

		sleepTime := targetDelta - time.Since(now).Seconds()
		if sleepTime > 0 {
			time.Sleep(time.Duration(sleepTime * float64(time.Second)))
		}

	}
}
