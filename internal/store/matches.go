package store

import (
	"context"
	"database/sql"
	"log"
)

type MatchStore struct {
	db *sql.DB
}

func (m *MatchStore) Create(ctx context.Context, match *Match) error {
	query := `
		INSERT INTO matches (player1_id, player2_id, winner_id, question_id)
        VALUES ($1, $2, $3, $4)
	`

	_, err := m.db.ExecContext(ctx, query, match.Player1ID, match.Player2ID, match.WinnerID, match.QuestionID)
	return err
}

func (m *MatchStore) GetMatchesWonByUser(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM matches 
		WHERE winner_id = $1
	`

	var count int
	err := m.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		log.Printf("Error fetching matches won for user %d: %v", userID, err)
		return 0, err
	}
	return count, nil
}

func (m *MatchStore) GetMatchesPlayedByUser(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM matches 
		WHERE player1_id = $1 OR player2_id = $1
	`

	var count int
	err := m.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		log.Printf("Error fetching matches played for user %d: %v", userID, err)
		return 0, err
	}
	return count, nil
}

func (m *MatchStore) GetTotalMatchesPlayed(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM matches
	`

	var count int
	err := m.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error fetching total matches played: %v", err)
		return 0, err
	}
	return count, nil
}
