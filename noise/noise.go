package noise

import (
	"math"
	"math/rand"
	"time"

	"github.com/ffardo/go-event-vision"
)

type RandomEventGenerator struct {
	random *rand.Rand
	factor float64
	maxX   int
	maxY   int
}

func (r RandomEventGenerator) generateRandomEvent(ts int) event.Event {
	x := r.random.Intn(r.maxX)
	y := r.random.Intn(r.maxY)

	pol := 0
	if rand.Intn(100) > 50 {
		pol = 1
	}

	return event.Event{
		Coords: event.Point2D{
			X: x, Y: y,
		},
		P:  pol,
		Ts: ts,
	}
}

func (r RandomEventGenerator) createAdditionalNoisyEvents(ts int) []event.Event {

	nE := make([]event.Event, 0)

	f := r.factor
	for f > 0.0 {
		lf := math.Min(1.0, f)
		if r.random.Float64() > 1.0-lf {
			nE = append(nE, r.generateRandomEvent(ts))
		}
		f -= 1.0
	}
	return nE
}

func (r RandomEventGenerator) flipCoin() bool {
	return r.random.Float64() > 1.0-r.factor
}

//ApplyAdditiveNoise inserts additive noise events into event stream
func ApplyAdditiveNoise(events []event.Event, maxX, maxY int, factor float64) []event.Event {
	noiseData := make([]event.Event, 0)
	r := RandomEventGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		factor: factor,
		maxX:   maxX,
		maxY:   maxY,
	}

	for _, ev := range events {

		noiseData = append(noiseData, ev)
		noiseData = append(noiseData, r.createAdditionalNoisyEvents(ev.Ts)...)
	}
	return noiseData
}

//ApplyDegenerativeNoise replaces some events with random events at the same timestamp
func ApplyDegenerativeNoise(events []event.Event, maxX, maxY int, factor float64) []event.Event {
	noiseData := make([]event.Event, len(events))
	r := RandomEventGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		factor: factor,
		maxX:   maxX,
		maxY:   maxY,
	}

	for i, ev := range events {
		if r.flipCoin() {
			noiseData[i] = r.generateRandomEvent(ev.Ts)
		} else {
			noiseData[i] = ev
		}
	}
	return noiseData
}
