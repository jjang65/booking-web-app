package dbrepo

import (
	"context"
	"errors"
	"github.com/jjang65/booking-web-app/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	log.Println("r.RestrictionID: ", r.RestrictionID)
	if err != nil {
		log.Println("InsertRoomRestriction::err:", err)
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns ture if availability exists for roomID, and false if no availability
func (m *postgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT COUNT(id)
				FROM room_restrictions
				WHERE room_id = $1
				    AND $2 <= end_date and $3 >= start_date;`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		log.Println("SearchAvailabilityByDatesByRoomID::err:", err)
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT r.id, r.room_name
			FROM rooms r
			WHERE r.id not in (
				SELECT room_id FROM room_restrictions rr 
					WHERE $1 < rr.end_date
						AND $2 > rr.start_date
			);
	`
	var rooms []models.Room
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDbRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		SELECT id, room_name, created_at, updated_at
			FROM rooms
			WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID returns a user by ID
func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
			FROM users
			WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser updates a user in the db
func (m *postgresDbRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE users SET first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5 
	`
	_, err := m.DB.ExecContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Authenticate authenticates a user
func (m *postgresDbRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `
		SELECT id, password
			FROM users
			WHERE email = $1
	`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations
func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, 
			r.created_at, r.updated_at, rm.id, rm.room_name
			FROM reservations r
			LEFT JOIN rooms rm ON (r.room_id = rm.id)
			ORDER BY r.start_date ASC
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// AllNewReservations returns a slice of all reservations
func (m *postgresDbRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, 
			r.created_at, r.updated_at, rm.id, rm.room_name, r.processed
			FROM reservations r
			LEFT JOIN rooms rm ON (r.room_id = rm.id)
			WHERE processed = 0
			ORDER BY r.start_date ASC
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
			&i.Processed,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}
