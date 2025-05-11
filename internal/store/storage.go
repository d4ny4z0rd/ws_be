package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *User) error
		Update(context.Context, int64, map[string]interface{}) error
		Delete(context.Context, int64) error
		GetByEmail(context.Context, string) (*User, error)
		IncrementPoints(context.Context, int64, int) (int, error)
		DecrementPoints(context.Context, int64, int) (int, error)
	}
	Matches interface {
		Create(context.Context, *Match) error
		GetMatchesWonByUser(context.Context, int64) (int, error)
		GetMatchesPlayedByUser(context.Context, int64) (int, error)
		GetTotalMatchesPlayed(context.Context) (int, error)
	}
	Questions interface {
		Create(context.Context, *DSAQuestion) error
		GetRandomQuestion(context.Context) (*DSAQuestion, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:     &UserStore{db},
		Matches:   &MatchStore{db},
		Questions: &QuestionStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
