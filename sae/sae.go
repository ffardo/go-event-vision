package sae

import (
	"errors"

	"github.com/ffardo/go-event-vision"
)

const (
	METHOD_ADDITIVE string = "additive"
	METHOD_RECENT   string = "recent"
)

func newMatrix(width, height int) ([][]int, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("Invalid matrix size")
	}
	m := make([][]int, height)

	for i := 0; i < height; i++ {
		m[i] = make([]int, width)
	}
	return m, nil
}

func createMatrix(events []event.Event, width, height int, method func([][]int, event.Event)) ([][]int, error) {
	m, err := newMatrix(width, height)
	if err != nil {
		return nil, err
	}
	for _, e := range events {
		if e.Coords.X >= 0 && e.Coords.X < width && e.Coords.Y >= 0 && e.Coords.Y < height {
			method(m, e)
		}
	}
	return m, nil
}

// CreateMatrix creates a Surface of Active Events in the form of a 2D matrix
// Supports nost recent event and accumulation by adding time stamps
func CreateMatrix(events []event.Event, method string, witdh, height int) ([][]int, error) {
	switch method {
	case METHOD_ADDITIVE:
		{
			return createMatrix(events, witdh, height, func(m [][]int, e event.Event) { m[e.Coords.Y][e.Coords.X] += e.Ts })
		}
	case METHOD_RECENT:
		{
			return createMatrix(events, witdh, height, func(m [][]int, e event.Event) { m[e.Coords.Y][e.Coords.X] = e.Ts })
		}
	}
	return nil, errors.New("Invalid method")
}

func createMap(events []event.Event, method func(map[event.Point2D]int, event.Event)) map[event.Point2D]int {
	m := make(map[event.Point2D]int)

	for _, e := range events {
		method(m, e)
	}
	return m
}

// CreateMap creates a Surface of Active Events in the form of a map.
// Supports nost recent event and accumulation by adding time stamps
func CreateMap(events []event.Event, method string) (map[event.Point2D]int, error) {
	switch method {
	case METHOD_ADDITIVE:
		{
			return createMap(events, func(m map[event.Point2D]int, e event.Event) { m[e.Coords] += e.Ts }), nil
		}
	case METHOD_RECENT:
		{
			return createMap(events, func(m map[event.Point2D]int, e event.Event) { m[e.Coords] = e.Ts }), nil
		}
	}
	return nil, errors.New("Invalid method")
}

// CropMap crops Time Surface map starting at x,y and ending at x+widht, y+height
func CropMap(src map[event.Point2D]int, x, y, width, height int) map[event.Point2D]int {
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
