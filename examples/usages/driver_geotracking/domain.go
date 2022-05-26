package main

// This package contains all domain entities for this sample.
// Consider it a `domain`/`entities`/`whatever you call it` package of your project
type (
	DriverID int
	EventID  int

	Driver struct {
		ID DriverID
	}

	Coordinate struct {
		Lat float32
		Lng float32
	}

	Location struct {
		Coordinate
		City    string
		State   string
		Country string
	}

	GeoDriver struct {
		Driver
		Location
	}
)
