# go-event-vision
Go package for event based vision

This project contains a port of some functions from [event-Python](https://github.com/gorchard/event-Python) by Garrick Orchard, and I wrote it to help on my Ph.D. thesis.

It also contains additional functionalities such as controlled noise applied directly to the event stream.

# Event Vision

Event-based cameras are imaging sensors that register local brightness changes in the form of asynchronous events. Each event is represented by its (x,y) coordinates, a timestamp, and a polarity indicating an upward or downward change. However, an event does not contain information about light intensity or magnitude.

In such sensors, every pixel is independent of the other pixels, and data is output as a stream of asynchronous events. The following table has an example of such event data.


| timestamp |  x | y | polarity |
|-----------|----|---|----------|
|    120    |  30| 40|     1    |
|    123    |  33| 43|     0    |
|    123    | 124| 80|     1    |
|    ...    | ...|...|    ...   |



This [video](https://www.youtube.com/watch?v=LauQ6LWTkxM) by the Robotics and Perception Group from the University of Zurich demonstrates how these sensors are different from conventional cameras with a global shutter.

# Why Go?

I initially developed my algorithms in Python and used [scikit-learn](https://scikit-learn.org/) for data cleanup and classification. However, processing raw data from N-MNIST and N-Caltech datasets in pure Python on CPU proved to be too time-consuming. With the added noise, the M-NIST dataset could get as large as 200+ GB.

I needed an alternative language that was fast and productive. I also would need to export processed data into a CSV file to apply machine learning scripts using Python at a later stage. Go seemed like a good alternative for the task and this change saved me precious hours.

Also, at the present moment, Go is not the de facto language for scientific research and I think new libraries and tools are always welcome.

# Features

This package features the following functionality

* ATIS format support
* Support to N-Caltech and N-MNIST datasets, including saccade stabilization
* Spatio-temporal filtering
* Refraction
* Additive and degenerative noise applied directly to the event stream

# Basic types

All (x,y) pixel coordinates are represented by the type ```Point2D```

```
type Point2D struct {
	X, Y int
}

```

An event is represented by the ```Event``` type which has the following format

```
type Event struct {
	Coords Point2D
	Ts     int
	P      int
}
```

All provided processing functions will accept a slice of ```Event``` as argument and return another slice with the processed events. No function will affect instance data. In this regard, go-event-vision differs slightly from [event-Python](https://github.com/gorchard/event-Python)

An event capture from a dataset is represented by the following ```EventCapture``` structure.

```
type EventCapture struct {
	Events []Event
	Width  int
	Height int
}
```

# Code example

The following code example shows the basic functionality of the event vision library.

```
package main

import (
	"log"

	"github.com/ffardo/go-event-vision/datasets"
	"github.com/ffardo/go-event-vision/datasets/neuromorphic"
	"github.com/ffardo/go-event-vision/filter"
)

func main() {

	//Adjust path to actual location
	reader := neuromorphic.NeuromorphicDataset{
		FilePath: "Caltech101/accordion/image_0005.bin",
	}

	evCap, err := datasets.ReadDataset(reader)

	if err != nil {
		log.Fatal(err)
	}

	//Appplying spatio-temporal filtering 5us
	evCap.Events = filter.SpatioTemporal(evCap.Events, evCap.Width, evCap.Height, 5000)

	//Applying refraction of 1us
	evCap.Events = filter.ApplyRefraction(evCap.Events, 1000)

	//Stabilize saccadic movements
	evCap.Events = neuromorphic.Stabilize(evCap.Events)

}
```

# Roadmap

This project is a work in progress and there is no tagged release yet. The following requirements and features are planned

* Full test coverage
* Instalation and documentation from pkg.go.dev
* Additional formats such as Prophesee
* Support to additional datasets such as N-Cars and DDD17


# Additional Information
Sample data in the tests contain a subset of events from a capture extracted from the N-Caltech dataset, which can be obtained  [here](https://www.garrickorchard.com/datasets/n-caltech101) [1]

# References

    [1] Orchard, G.; Cohen, G.; Jayawant, A.; and Thakor, N.  “Converting Static Image Datasets to Spiking Neuromorphic Datasets Using Saccades", Frontiers in Neuroscience, vol.9, no.437, Oct. 2015 (open access Frontiers link)