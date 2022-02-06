package ncars

import (
	"github.com/ffardo/go-event-vision"
	"github.com/ffardo/go-event-vision/format/prophesee"
)

// Ncars implements DatasetReader interface for N-CARS dataset
type Ncars struct {
	FilePath string
}

// Read event capture for an entry in the dataset
func (n Ncars) Read() (event.EventCapture, error) {
	atis := prophesee.Dat{FilePath: n.FilePath}

	return atis.ReadEvents()
}

// Write capture to a dataset. Should be used only for data augmentation.
func (n Ncars) Write(evCap event.EventCapture) error {
	atis := prophesee.Dat{FilePath: n.FilePath}

	return atis.WriteEvents(evCap)
}
