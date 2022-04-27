package dbrepo

import (
	"github.com/jjang65/booking-web-app/internal/models"
	"time"
)

func (m *testDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into db that returns reservation_id and error
func (m *testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into db
func (m *testDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID returns ture if availability exists for roomID, and false if no availability
func (m *testDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDbRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}
