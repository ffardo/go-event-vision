package render

import (
	"image"
	"image/color"

	"github.com/ffardo/go-event-vision"
)

// SaeMap renders normalized SAE data in map format to an image.RGBA pointer
func SaeMap(sae map[event.Point2D]int, width, height int) *image.RGBA {
	image := image.NewRGBA(image.Rectangle{image.Pt(0, 0), image.Pt(width, height)})

	bg := color.RGBA{A: 255}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			image.Set(j, i, bg)
		}
	}

	max := 0
	for _, v := range sae {
		if v > max {
			max = v
		}
	}

	maxF := float64(max)

	for pt, v := range sae {
		c := uint8((float64(v) / maxF) * 255.0)
		image.Set(pt.X, pt.Y, color.RGBA{A: 255, R: c, B: c, G: c})
	}
	return image
}
