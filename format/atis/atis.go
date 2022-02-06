package atis

import (
	"math"
	"os"

	"github.com/ffardo/go-event-vision"
)

// Aer implements ATIS AER format reading and writing
type Aer struct {
	FilePath string
}

func (a Aer) newEventFromBytes(data []byte) event.Event {
	x := int(data[0])
	y := int(data[1])
	ts1 := (int(data[2]) & 127) << 16
	ts2 := int(data[3]) << 8
	ts3 := int(data[4])
	ts := ts1 + ts2 + ts3
	p := (int(data[2]) & 128) >> 7

	if y == 240 {
		ts += int(math.Pow(2, 13))
	}

	return event.Event{
		Coords: event.Point2D{X: x, Y: y},
		P:      p,
		Ts:     ts,
	}
}

func (a Aer) eventToBytes(ev event.Event) []byte {
	x := ev.Coords.X
	y := ev.Coords.Y

	ts := ev.Ts

	if y == 240 {
		ts -= int(math.Pow(2, 13))
	}

	data := make([]byte, 5)

	data[0] = byte(x)
	data[1] = byte(y)

	p := ev.P << 7
	data[2] = byte(((ts >> 16) & 127) + p)
	data[3] = byte((ts >> 8) & 255)
	data[4] = byte(ts & 255)

	return data
}

// ReadEvents read events in the ATIS AER format from file
func (a Aer) ReadEvents() (event.EventCapture, error) {
	f, err := os.Open(a.FilePath)
	if err != nil {
		return event.EventCapture{}, err
	}

	defer f.Close()

	ev := []event.Event{}

	bb := make([]byte, 5)
	bc := 5

	mX, mY := 0, 0

	for bc == 5 {
		bc, err = f.Read(bb)
		if bc == 5 && err == nil {
			n := a.newEventFromBytes(bb)
			if n.Coords.Y == 240 {
				continue
			}
			if n.Coords.X > mX {
				mX = n.Coords.X
			}
			if n.Coords.Y > mY {
				mY = n.Coords.Y
			}
			ev = append(ev, n)
		}
	}

	return event.EventCapture{
		Events: ev,
		Width:  mX + 1,
		Height: mY + 1,
	}, nil
}

// WriteEvents will write events to file in the ATIS AER format
func (a Aer) WriteEvents(evCap event.EventCapture) error {
	f, err := os.Create(a.FilePath)

	defer f.Close()
	if err != nil {
		return err
	}

	for _, ev := range evCap.Events {

		data := a.eventToBytes(ev)
		f.Write(data)
	}
	return nil
}
