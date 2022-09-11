package render

import (
	"image"
	"image/color"

	"github.com/ffardo/go-event-vision"
)

// Stream renders an event stream to an image.RGBA pointer.
func Stream(events []event.Event, width, height int, background, positive, negative color.RGBA) *image.RGBA {
	image := image.NewRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(width, height)})

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			image.Set(j, i, background)
		}
	}

	for _, e := range events {
		if e.P == 1 {
			image.Set(e.Coords.X, e.Coords.Y, positive)
		} else {
			image.Set(e.Coords.X, e.Coords.Y, negative)
		}
	}

	return image
}
