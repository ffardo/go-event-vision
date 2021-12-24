package noise

import (
	"math"
	"math/rand"
	"time"

	"github.com/ffardo/go-event-vision"
)

type randomEventGenerator struct {
	random *rand.Rand
	factor float64
	maxX   int
	maxY   int
}

func (r randomEventGenerator) generateRandomEvent(ts int) event.Event {
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

func (r randomEventGenerator) createAdditionalNoisyEvents(ts int) []event.Event {
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

func (r randomEventGenerator) flipCoin() bool {
	return r.random.Float64() > 1.0-r.factor
}

// ApplyAdditive inserts additive noise events into event stream.  Since new events are added to the stream, applying additive noise will
// result in a larger event slice.
// The factor argument will determine the likelihood of a random event to be added to the stream for each existing event.
// If factor is set to 0.0, no new events will be added.
// If factor is set to 0.25 there is a 25% chance that a new event will be added to the stream at the same time as each event.
// If factor is set to 5.25, it is certain that 5 random events will be added at the same timestamp for each event in the stream and
// there is a 25% chance of a sixth random event to be added t the stream.
func ApplyAdditive(src []event.Event, maxX, maxY int, factor float64) []event.Event {
	dst := make([]event.Event, 0)
	r := randomEventGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		factor: factor,
		maxX:   maxX,
		maxY:   maxY,
	}

	for _, ev := range src {

		dst = append(dst, ev)
		dst = append(dst, r.createAdditionalNoisyEvents(ev.Ts)...)
	}
	return dst
}

// ApplyDegenerative replaces some events in en event stream with random events at the same timestamp.
// Since no new event is added, the resulting slice will have the same size as the source event slice.
// The amount of affected events will vary depending on the factor value.
// If factor is set to 0.0 no event will be affected.
// If factor is set to 0.25, 25% of the events are likely to be affected.
// If factor is set to 1.0, all events will be replaced by random events at the same timestamps.
func ApplyDegenerative(src []event.Event, maxX, maxY int, factor float64) []event.Event {
	dst := make([]event.Event, len(src))
	r := randomEventGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		factor: factor,
		maxX:   maxX,
		maxY:   maxY,
	}

	for i, ev := range src {
		if r.flipCoin() {
			dst[i] = r.generateRandomEvent(ev.Ts)
		} else {
			dst[i] = ev
		}
	}
	return dst
}
