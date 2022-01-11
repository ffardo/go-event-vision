// event is the base package for go-event vision and contains base struct definitions
package event

// Point2D represents the 2D coordinates of an event
type Point2D struct {
	X, Y int
}

// Event represents a discrete event with coordinates (X,Y), timestamp (Ts) and polarity (P)
type Event struct {
	Coords Point2D
	Ts     int
	P      int
}

// EventCapture represents a scene captured with event sensors
type EventCapture struct {
	Events []Event
	Width  int
	Height int
}
