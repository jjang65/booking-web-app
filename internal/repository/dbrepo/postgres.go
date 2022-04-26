package dbrepo

import (
	"context"
	"github.com/jjang65/booking-web-app/internal/models"
	"log"
	"time"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into db that returns reservation_id and error
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, 
			end_date, room_id, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	// QueryRowContext executes a statement and returns the most recent row
	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into db
func (m *postgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id,
			created_at, updated_at, restriction_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		log.Println("InsertRoomRestriction::err:", err)
		return err
	}
	return nil
}

// SearchAvailabilityByDates returns ture if availability exists for roomID, and false if no availability
func (m *postgresDbRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT COUNT(id)
				FROM room_restrictions
				WHERE room_id $1
				    AND $2 < end_date and $3 > start_date;`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		log.Println("SearchAvailabilityByDates::err:", err)
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}
