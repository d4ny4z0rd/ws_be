package store

import (
	"context"
	"database/sql"
)

type QuestionStore struct {
	db *sql.DB
}

func (s *QuestionStore) create(ctx context.Context, tx *sql.Tx, q *DSAQuestion) error {
	query := `
		INSERT INTO dsa_questions 
		(title, description, input_format, output_format, example_input, example_output)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		q.Title,
		q.Description,
		q.InputFormat,
		q.OutputFormat,
		q.ExampleInput,
		q.ExampleOutput,
	).Scan(&q.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *QuestionStore) Create(ctx context.Context, q *DSAQuestion) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, q)
	})
}

func (s *QuestionStore) GetRandomQuestion(ctx context.Context) (*DSAQuestion, error) {

	query := `
		SELECT id, title, description, input_format, output_format, example_input, example_output
		FROM dsa_questions
		ORDER BY RANDOM()
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var q DSAQuestion
	err := s.db.QueryRowContext(ctx, query).Scan(
		&q.ID,
		&q.Title,
		&q.Description,
		&q.InputFormat,
		&q.OutputFormat,
		&q.ExampleInput,
		&q.ExampleOutput,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &q, nil
}
