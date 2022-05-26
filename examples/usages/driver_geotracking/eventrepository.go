package main

import (
	"context"
	"fmt"
)

type (
	// EventRepository is a repository for the event resource that allows us
	// to unarmshall it to get its inner data or flag it as processed to avoid
	// processing it many times.
	EventRepository struct {
		// stuff to get event information (eg. SQL / http.Client / etc)
	}
)

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) GetTrigger(ctx context.Context, e EventID) (DriverID, error) {
	// use stuff to get driver id from event
	fmt.Printf("getting driver id by event id %d\n", e)
	return DriverID(1234), nil
}

func (r *EventRepository) MarkProcessed(ctx context.Context, id EventID) error {
	// use stuff to mark the event as processed
	fmt.Printf("marking event id %d as processed\n", id)
	return nil
}
