package neuromorphic

import (
	"reflect"
	"testing"

	"github.com/ffardo/go-event-vision"
)

func TestNeuromorphicDataset_Read(t *testing.T) {

	refData := event.EventCapture{
		Events: []event.Event{
			{Coords: event.Point2D{X: 10, Y: 30}, Ts: 937, P: 1},
			{Coords: event.Point2D{X: 33, Y: 20}, Ts: 1030, P: 1},
			{Coords: event.Point2D{X: 12, Y: 27}, Ts: 1052, P: 1},
			{Coords: event.Point2D{X: 33, Y: 3}, Ts: 2078, P: 1},
			{Coords: event.Point2D{X: 14, Y: 23}, Ts: 2383, P: 0},
			{Coords: event.Point2D{X: 16, Y: 10}, Ts: 3189, P: 0},
			{Coords: event.Point2D{X: 7, Y: 30}, Ts: 4003, P: 1},
			{Coords: event.Point2D{X: 1, Y: 28}, Ts: 4975, P: 1},
			{Coords: event.Point2D{X: 16, Y: 11}, Ts: 6609, P: 0},
			{Coords: event.Point2D{X: 18, Y: 19}, Ts: 6678, P: 1},
		},
		Width:  34,
		Height: 31,
	}

	tests := []struct {
		name    string
		n       NeuromorphicDataset
		want    event.EventCapture
		wantErr bool
	}{
		{name: "Read valid dataset", n: NeuromorphicDataset{FilePath: "../../sample_data/neuro_sample.bin"}, want: refData, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Read()

			if (err != nil) != tt.wantErr {
				t.Errorf("NeuromorphicDataset.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NeuromorphicDataset.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
