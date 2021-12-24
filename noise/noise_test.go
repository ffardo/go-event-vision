package noise

import (
	"reflect"
	"testing"

	"github.com/ffardo/go-event-vision"
)

func TestApplyAdditive(t *testing.T) {

	src := []event.Event{
		{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
	}

	if !reflect.DeepEqual([]event.Event{}, ApplyAdditive([]event.Event{}, 10, 10, 1.0)) {
		t.Errorf("When source slice is empty, result slice should also be empty")
	}

	noised0 := ApplyAdditive(src, 10, 10, 0.0)

	if !reflect.DeepEqual(src, noised0) {
		t.Errorf("When factor is 0.0, result slice should be equal to source slice ")
	}

	noised1 := ApplyAdditive(src, 10, 10, 1.0)

	if len(noised1) != 2 {
		t.Errorf("when factor is set to 1.0 resulting slice should have an additional event")
	}

	if reflect.DeepEqual(src, noised1) {
		t.Errorf("When factor is set to 1.0 resulting slice should be different to source slice")
	}

}

func TestApplyDegenerative(t *testing.T) {

	src := []event.Event{
		{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
		{Coords: event.Point2D{X: 2, Y: 2}, Ts: 1, P: 1},
		{Coords: event.Point2D{X: 3, Y: 3}, Ts: 1, P: 1},
		{Coords: event.Point2D{X: 4, Y: 4}, Ts: 1, P: 1},
		{Coords: event.Point2D{X: 5, Y: 5}, Ts: 1, P: 1},
	}

	if !reflect.DeepEqual([]event.Event{}, ApplyDegenerative([]event.Event{}, 10, 10, 1.0)) {
		t.Errorf("When source slice is empty, result slice should also be empty")
	}

	noised0 := ApplyDegenerative(src, 10, 10, 0.0)

	if !reflect.DeepEqual(src, noised0) {
		t.Errorf("When factor is 0.0, result slice should be equal to source slice ")
	}

	noised1 := ApplyDegenerative(src, 10, 10, 1.0)

	if len(noised1) != len(src) {
		t.Errorf("Result slice should always have the same length as source slice")
	}

	if reflect.DeepEqual(src, noised1) {
		t.Errorf("When factor is set to 1.0 resulting slice should be different to source slice")
	}

}
