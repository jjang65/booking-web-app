package repository

import "github.com/jjang65/booking-web-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}
