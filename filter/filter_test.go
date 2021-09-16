package filter

import (
	"reflect"
	"testing"

	"github.com/ffardo/go-event-vision"
)

func TestCropTsMap(t *testing.T) {
	type args struct {
		src    map[event.Point2D]int
		x      int
		y      int
		width  int
		height int
	}

	tsMap := map[event.Point2D]int{
		{X: 1, Y: 1}: 1,
		{X: 5, Y: 5}: 2,
	}

	tests := []struct {
		name string
		args args
		want map[event.Point2D]int
	}{
		{name: "Test simple crop with one event inside area", args: args{tsMap, 0, 0, 2, 2}, want: map[event.Point2D]int{{X: 1, Y: 1}: 1}},
		{name: "Test crop area between events", args: args{tsMap, 2, 2, 3, 3}, want: map[event.Point2D]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CropTsMap(tt.args.src, tt.args.x, tt.args.y, tt.args.width, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CropTsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpatioTemporal(t *testing.T) {
	type args struct {
		src    []event.Event
		xMax   int
		yMax   int
		usTime int
	}

	srcData := []event.Event{
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
	}

	tests := []struct {
		name string
		args args
		want []event.Event
	}{
		{name: "Test sample event filtering",
			args: args{src: srcData, xMax: 35, yMax: 35, usTime: 5000},
			want: []event.Event{
				{Coords: event.Point2D{X: 16, Y: 11}, Ts: 6609, P: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpatioTemporal(tt.args.src, tt.args.xMax, tt.args.yMax, tt.args.usTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpatioTemporal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyRefraction(t *testing.T) {
	type args struct {
		src    []event.Event
		usTime int
	}

	tests := []struct {
		name string
		args args
		want []event.Event
	}{
		{
			name: "Test refraction two events at same location within interval",
			args: args{
				src: []event.Event{
					{Coords: event.Point2D{X: 10, Y: 30}, Ts: 900, P: 1},
					{Coords: event.Point2D{X: 10, Y: 30}, Ts: 1000, P: 1},
				},
				usTime: 1000},
			want: []event.Event{
				{Coords: event.Point2D{X: 10, Y: 30}, Ts: 900, P: 1},
			},
		},
		{
			name: "Test refraction two events at different location within interval",
			args: args{
				src: []event.Event{
					{Coords: event.Point2D{X: 10, Y: 30}, Ts: 900, P: 1},
					{Coords: event.Point2D{X: 11, Y: 30}, Ts: 1000, P: 1},
				},
				usTime: 1000},
			want: []event.Event{
				{Coords: event.Point2D{X: 10, Y: 30}, Ts: 900, P: 1},
				{Coords: event.Point2D{X: 11, Y: 30}, Ts: 1000, P: 1},
			},
		},
		{
			name: "Test refraction two events at same location outside interval",
			args: args{
				src: []event.Event{
					{Coords: event.Point2D{X: 10, Y: 30}, Ts: 1000, P: 1},
					{Coords: event.Point2D{X: 10, Y: 30}, Ts: 3000, P: 1},
				},
				usTime: 1000},
			want: []event.Event{
				{Coords: event.Point2D{X: 10, Y: 30}, Ts: 1000, P: 1},
				{Coords: event.Point2D{X: 10, Y: 30}, Ts: 3000, P: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyRefraction(tt.args.src, tt.args.usTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyRefraction() = %v, want %v", got, tt.want)
			}
		})
	}
}
