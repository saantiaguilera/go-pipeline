package main

import (
	"context"
	"fmt"
)

type (
	// TrackingRepository allows us to determine if a driver is close or not to a specific destination
	TrackingRepository struct {
		// stuff to get tracking info (eg. SQL / http.Client / etc)
	}
)

func NewTrackingRepository() *TrackingRepository {
	return &TrackingRepository{}
}

func (r *TrackingRepository) IsDriverClose(ctx context.Context, d GeoDriver) bool {
	// use stuff to determine if driver is close to a destination
	fmt.Printf("checking if driver is close to destination %+v\n", d)
	fmt.Println("driver is close to destination!")
	return true
}
