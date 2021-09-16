package filter

import (
	"math"

	"github.com/ffardo/go-event-vision"
)

func intMax(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//CropTsMap crops Time Surface map starting at x,y and ending at x+widht, y+height
func CropTsMap(src map[event.Point2D]int, x, y, width, height int) map[event.Point2D]int {
	dst := make(map[event.Point2D]int)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			c := event.Point2D{X: j + x, Y: i + y}
			value, ok := src[c]
			if ok {
				dst[c] = value
			}
		}
	}

	return dst
}

/*
SpatioTemporal generate a filtered set of events.
Uses a background activity filter on the events, such that only events which are
correlated with a neighbouring event within 'usTime' microseconds will be allowed
through the filter.
*/
func SpatioTemporal(src []event.Event, xMax, yMax, usTime int) []event.Event {

	t0 := make(map[event.Point2D]int)

	for _, ev := range src {
		t0[ev.Coords] = -usTime
	}

	xP := 0
	yP := 0
	pP := 0

	ex := make([]bool, len(src))

	totalEvents := len(src)

	for idx, dt := range src {
		if xP != dt.Coords.X || yP != dt.Coords.Y || pP != dt.P {
			t0[dt.Coords] = -usTime
			minXSub := intMax(0, dt.Coords.X-1)
			maxXSub := intMin(xMax, dt.Coords.X+1)
			minYSub := intMax(0, dt.Coords.Y-1)
			maxYSub := intMin(yMax, dt.Coords.Y+1)

			t0Temp := CropTsMap(t0, minXSub, minYSub, (maxXSub-minXSub)+1, (maxYSub-minYSub)+1)

			minTs := int(math.MaxInt64)
			for _, v := range t0Temp {
				minTs = intMin(minTs, dt.Ts-v)
			}

			if minTs > usTime {
				ex[idx] = true
				totalEvents--
			}
		}
		t0[dt.Coords] = dt.Ts
		xP = dt.Coords.X
		yP = dt.Coords.Y
		pP = dt.P

	}

	dst := make([]event.Event, totalEvents)
	pos := 0

	for idx, v := range ex {
		if v != true {
			dst[pos] = src[idx]
			pos++
		}
	}

	return dst
}

/*
ApplyRefraction applies a refractory period for each event.
In other words, if an event occurs within 'usTime' microseconds of a previous event at the
same pixel, then the second event is removed
*/
func ApplyRefraction(src []event.Event, usTime int) []event.Event {
	t0 := make(map[event.Point2D]int)

	for _, e := range src {
		t0[e.Coords] = -usTime
	}

	ex := make([]bool, len(src))

	totalEvents := len(src)

	for idx, dt := range src {
		if dt.Ts-t0[dt.Coords] < usTime {
			ex[idx] = true
			totalEvents--
		} else {
			t0[dt.Coords] = dt.Ts
		}
	}

	dst := make([]event.Event, totalEvents)
	pos := 0
	for idx, v := range ex {
		if v != true {
			dst[pos] = src[idx]
			pos++
		}
	}

	return dst
}
