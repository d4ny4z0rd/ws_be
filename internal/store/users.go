package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with the email already exists")
	ErrDuplicateUsername = errors.New("a user with the username already exists")
)

type password struct {
	text *string
	hash []byte
}

type UserStore struct {
	db *sql.DB
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) GetHash() []byte {
	return p.hash	
}

func (s *UserStore) create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (email, password, username) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Username,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq : duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, user)
	})
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, email, password, username, points, created_at, updated_at
		FROM users
		WHERE users.id = $1
	`

	user := &User{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, userID int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("no new fields to update")
	}

	query := `UPDATE users SET`
	args := []interface{}{}
	counter := 1

	for key, val := range updates {
		query += fmt.Sprintf(" %s = $%d,", key, counter)
		args = append(args, val)
		counter++
	}

	query = strings.TrimSuffix(query, ",") + fmt.Sprintf(" WHERE id = $%d", counter)
	args = append(args, userID)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (s *UserStore) Update(ctx context.Context, userID int64, updates map[string]interface{}) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.update(ctx, tx, userID, updates)
	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.delete(ctx, tx, userID)
	})
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, username, password, points, created_at, updated_at FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password, // Password will be directly scanned as []byte
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m *UserStore) IncrementPoints(ctx context.Context, userID int64, amount int) (int, error) {
	query := `UPDATE users SET points = points + $1 WHERE id = $2 RETURNING points`
	var newPoints int
	err := m.db.QueryRowContext(ctx, query, amount, userID).Scan(&newPoints)
	return newPoints, err
}

func (m *UserStore) DecrementPoints(ctx context.Context, userID int64, amount int) (int, error) {
	query := `UPDATE users SET points = GREATEST(points - $1, 0) WHERE id = $2 RETURNING points`
	var newPoints int
	err := m.db.QueryRowContext(ctx, query, amount, userID).Scan(&newPoints)
	return newPoints, err
}
