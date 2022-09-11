package prophesee

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/ffardo/go-event-vision"
)

// Dat implements Prophesee RAW DAT format reading and writing
// More information can be found in the official documentation
// https://docs.prophesee.ai/stable/data_formats/file_formats/dat.html
type Dat struct {
	FilePath string
}

func (d Dat) newEventFromBytes(data []byte) event.Event {

	ts := int(binary.LittleEndian.Uint32(data[:4]))
	addressValue := int(binary.LittleEndian.Uint32(data[4:]))

	x := int((addressValue & 0x00003FFF))
	y := int((addressValue & 0x0FFFC000) >> 14)
	p := int((addressValue & 0x10000000) >> 28)

	return event.Event{Coords: event.Point2D{X: x, Y: y}, P: p, Ts: ts}
}

// ReadEvents read events in the Prophesee RAW DAT format from file
func (d Dat) ReadEvents() (event.EventCapture, error) {
	f, err := os.Open(d.FilePath)

	if err != nil {
		return event.EventCapture{}, err
	}

	defer f.Close()

	ev := []event.Event{}
	bc := 2

	mX, mY := 0, 0
	evSize := 0

	evSize, err = d.seekFirstEvent(f)
	bc = evSize

	bb := make([]byte, 8)
	for bc == evSize {
		bc, err = f.Read(bb)
		if bc == evSize && err == nil {

			n := d.newEventFromBytes(bb)
			if n.Coords.X > mX {
				mX = n.Coords.X
			}
			if n.Coords.Y > mY {
				mY = n.Coords.Y
			}
			ev = append(ev, n)
		}
	}

	return event.EventCapture{Events: ev, Width: mX + 1, Height: mY + 1}, nil
}

func (d Dat) seekFirstEvent(f *os.File) (int, error) {
	var err error
	bh := make([]byte, 2)
	bc := 2

	evSize := 0

	for bc == 2 {
		bc, err = f.Read(bh)
		if bh[0] == '%' && bh[1] == ' ' {
			bhi := make([]byte, 1)
			for bhi[0] != '\n' {
				_, err = f.Read(bhi)
			}
			continue
		}

		evSize = int(bh[1])
		break
	}

	return evSize, err
}

// WriteEvents will write events to file in the Prophesee RAW DAT format
func (d Dat) WriteEvents(evCap event.EventCapture) error {
	return errors.New("Prophesee.WriteEvents is not implemented")
}
