package store

import "time"

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  []byte `json:"-"`
	Username  string `json:"username"`
	Points    int    `json:"points"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DSAQuestion struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	InputFormat   string `json:"input_format"`
	OutputFormat  string `json:"output_format"`
	ExampleInput  string `json:"example_input"`
	ExampleOutput string `json:"example_output"`
}

type Match struct {
	ID         int64
	Player1ID  int64
	Player2ID  int64
	WinnerID   int64
	QuestionID int64
	CreatedAt  time.Time
}
