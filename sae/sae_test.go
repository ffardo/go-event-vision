package sae

import (
	"reflect"
	"testing"

	"github.com/ffardo/go-event-vision"
)

func TestCreateMatrix(t *testing.T) {
	type args struct {
		events []event.Event
		method string
		width  int
		height int
	}

	tests := []struct {
		name    string
		args    args
		want    [][]int
		wantErr bool
	}{
		{
			"Test invalid method",
			args{events: []event.Event{}, method: "unknown_method"},
			nil, true,
		},
		{
			"Test create matrix sae with zero width and height",
			args{
				events: []event.Event{},
				method: "additive",
				width:  0,
				height: 0,
			},
			nil, true,
		},
		{
			"Test create matrix sae with negative width and height",
			args{
				events: []event.Event{},
				method: "additive",
				width:  -1,
				height: -1,
			},
			nil, true,
		},
		{
			"Test with additive method and two events in the same location",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "additive",
				width:  2,
				height: 2,
			},
			[][]int{{0, 0}, {0, 3}},
			false,
		},
		{
			"Test with additive method and two events at different locations",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 0, Y: 0}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "additive",
				width:  2,
				height: 2,
			},
			[][]int{{1, 0}, {0, 2}},
			false,
		},
		{
			"Test with additive method and events outside matrix area",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 3, Y: 3}, Ts: 2, P: 1},
					{Coords: event.Point2D{X: -3, Y: -3}, Ts: 2, P: 1},
				},
				method: "additive",
				width:  2,
				height: 2,
			},
			[][]int{{0, 0}, {0, 0}},
			false,
		},
		{
			"Test with most recent method and two events in the same location",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "recent",
				width:  2,
				height: 2,
			},
			[][]int{{0, 0}, {0, 2}},
			false,
		},
		{
			"Test with most recent method and two events at diffent locations",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 0, Y: 0}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "recent",
				width:  2,
				height: 2,
			},
			[][]int{{1, 0}, {0, 2}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateMatrix(tt.args.events, tt.args.method, tt.args.width, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMatrixSae() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMatrixSae() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMap(t *testing.T) {
	type args struct {
		events []event.Event
		method string
	}

	tests := []struct {
		name    string
		args    args
		want    map[event.Point2D]int
		wantErr bool
	}{
		{
			"Test invalid method",
			args{
				events: []event.Event{},
				method: "unknown_method",
			},
			nil,
			true,
		},
		{
			"Test with additive method and two events in the same location",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "additive",
			},
			map[event.Point2D]int{
				{X: 1, Y: 1}: 3,
			},
			false,
		},
		{
			"Test with most recent method and two events in the same location",
			args{
				events: []event.Event{
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 1, P: 1},
					{Coords: event.Point2D{X: 1, Y: 1}, Ts: 2, P: 1},
				},
				method: "recent",
			},
			map[event.Point2D]int{
				{X: 1, Y: 1}: 2,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateMap(tt.args.events, tt.args.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMapSae() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMapSae() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCropMap(t *testing.T) {
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
			if got := CropMap(tt.args.src, tt.args.x, tt.args.y, tt.args.width, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CropMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
