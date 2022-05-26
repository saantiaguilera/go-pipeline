package main

import (
	"context"
	"fmt"
)

type (
	// DriverRepository is a repository that lets us interact with the
	// driver resource (eg. getting information, checking where it is in real time,
	// saving if it got to the destination, etc.)
	DriverRepository struct {
		// stuff to get a driver (eg. SQL / http.Client / etc)
	}
)

func NewDriverRepository() *DriverRepository {
	return &DriverRepository{}
}

func (r *DriverRepository) GetDriverByID(ctx context.Context, id DriverID) (Driver, error) {
	// use stuff to get driver by id
	fmt.Printf("getting driver data by id %d\n", id)
	return Driver{
		ID: id,
	}, nil
}

func (r *DriverRepository) GetRealTimeCoordinatesByID(ctx context.Context, id DriverID) (Coordinate, error) {
	// use stuff to get current coordinates of driver
	fmt.Printf("getting driver real time coordinates by id %d\n", id)
	return Coordinate{
		Lat: 32,
		Lng: 1,
	}, nil
}

func (r *DriverRepository) SaveDriverInDestination(ctx context.Context, d GeoDriver) error {
	// use stuff to save data
	fmt.Printf("saving driver in destination %+v\n", d)
	return nil
}
