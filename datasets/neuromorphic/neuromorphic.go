package neuromorphic

import (
	"github.com/ffardo/go-event-vision"
	"github.com/ffardo/go-event-vision/format"
)

//NeuromorphicDataset implements DatasetReader interface for N-MNIST and N-Caltech100 datasets
type NeuromorphicDataset struct {
	FilePath string
}

//Read event capture for an entry in the dataset
func (n NeuromorphicDataset) Read() (event.EventCapture, error) {

	atis := format.Atis{FilePath: n.FilePath}

	return atis.ReadEvents()

}

//Write capture to a dataset. Should be used only for data augmentation.
func (n NeuromorphicDataset) Write(evCap event.EventCapture) error {

	atis := format.Atis{FilePath: n.FilePath}

	return atis.WriteEvents(evCap)

}

// Stabilize saccadic motion
func Stabilize(src []event.Event) []event.Event {

	dst := []event.Event{}

	ts := 0

	s1, s2 := 0, 0

	for i := 0; i < len(src) && ts <= 105e3; i++ {
		ts = src[i].Ts
		s1 = i
	}

	for i := s1; i < len(src) && ts <= 210e3; i++ {
		ts = src[i].Ts
		s2 = i
	}

	dst = append(dst, correctSaccade1(src[0:s1])...)
	dst = append(dst, correctSaccade2(src[s1:s2])...)
	dst = append(dst, correctSaccade3(src[s2:])...)

	return dst
}

func correctSaccade1(ev []event.Event) []event.Event {

	c := make([]event.Event, len(ev))

	for i, ev := range ev {
		e := ev

		e.Coords.X -= int(3.5 * float64(e.Ts) / float64(105e3))
		e.Coords.Y -= int(3.5 * float64(e.Ts) / float64(105e3))
		c[i] = e
	}
	return c
}

func correctSaccade2(ev []event.Event) []event.Event {

	c := make([]event.Event, len(ev))

	for i, ev := range ev {
		e := ev

		e.Coords.X -= int(3.5 + 3.5*(float64(e.Ts)-105e3)/105e3)
		e.Coords.Y -= int(7 + 7*(float64(e.Ts)-105e3)/105e3)
		c[i] = e
	}
	return c
}

func correctSaccade3(ev []event.Event) []event.Event {

	c := make([]event.Event, len(ev))

	for i, ev := range ev {
		e := ev

		e.Coords.X -= int(7 + 7*(float64(e.Ts)-210e3)/105e3)
		c[i] = e
	}
	return c
}
