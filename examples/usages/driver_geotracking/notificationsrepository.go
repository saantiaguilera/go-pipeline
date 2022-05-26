package main

import (
	"context"
	"fmt"
)

type (
	// NotificationRepository allows us to send notifications to a driver
	NotificationRepository struct {
		// stuff to send notifications (http client / firebase / etc)
	}
)

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) NotifyCloseToDestination(ctx context.Context, d Driver) error {
	// send notification using struct stuff
	fmt.Printf("sending notification to driver %+v\n", d)
	return nil
}
