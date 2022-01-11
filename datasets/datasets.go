package datasets

import "github.com/ffardo/go-event-vision"

// DatasetReader specifies an interface to read dataset captures
type DatasetReader interface {
	Read() (event.EventCapture, error)
}

// DatasetWriter specifies an interface to write dataset captures
type DatasetWriter interface {
	Write(event.EventCapture) error
}

// ReadDataset reads event data from a DatasetReader interface
func ReadDataset(reader DatasetReader) (event.EventCapture, error) {
	return reader.Read()
}

// ReadDataset writes event data to a DatasetWriter interface. Should only be used for data augmentation
func WriteDataset(writer DatasetWriter, cap event.EventCapture) error {
	return writer.Write(cap)
}
