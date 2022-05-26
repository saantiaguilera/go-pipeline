package main

import (
	"context"
	"fmt"
)

type (
	// LocationRepository is a repository for the location resource that allows us
	// to get a full location from a given input of coordinates
	LocationRepository struct {
		// stuff to get location info (eg. SQL / http.Client / etc)
	}
)

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

func (r *LocationRepository) GetFullLocationByLatLng(ctx context.Context, c Coordinate) (Location, error) {
	// use stuff to get location
	fmt.Printf("getting full location by coordinates %+v\n", c)
	return Location{
		Coordinate: c,
		City:       "example",
		State:      "example",
		Country:    "example",
	}, nil
}
