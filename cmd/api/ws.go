package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"ws_practice_1/internal/env"
	"ws_practice_1/internal/store"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Match struct {
	Player1     *websocket.Conn
	Player2     *websocket.Conn
	Question    store.DSAQuestion
	IsCompleted bool
	mu          sync.Mutex
}

type wsApp struct {
	matchMaking MatchMaking
	matches     map[*websocket.Conn]*Match
	scores      map[*websocket.Conn]int
	userConns   map[int64]*websocket.Conn
	connUsers   map[*websocket.Conn]int64
	mu          sync.Mutex
	app         *application
	userData    map[int64]*store.User
}

type MatchMaking struct {
	mu      sync.Mutex
	waiting *websocket.Conn
}

type response struct {
	Type     string      `json:"type"`
	Message  interface{} `json:"message"`
	Opponent *Opponent   `json:"opponent,omitempty"`
}

type Opponent struct {
	Username string `json:"username"`
	Points   int    `json:"points"`
}

type payload struct {
	Type   string `json:"type"`
	Answer string `json:"answer"`
	LangID int    `json:"language_id"`
}

const judge0APIURL = "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=false&wait=true"

type submissionRequest struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin"`
}

type submissionResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Message  string `json:"message"`
	StatusID int    `json:"status_id"`
}

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	log.Println("User:", user)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	println("hello")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Printf("WebSocket connection established for user %d\n", user.ID)

	app.ws.mu.Lock()

	if app.ws.userConns == nil {
		app.ws.userConns = make(map[int64]*websocket.Conn)
	}
	if app.ws.connUsers == nil {
		app.ws.connUsers = make(map[*websocket.Conn]int64)
	}
	if app.ws.userData == nil {
		app.ws.userData = make(map[int64]*store.User)
	}

	if oldConn, exists := app.ws.userConns[user.ID]; exists {
		if app.ws.matchMaking.waiting == oldConn {
			app.ws.matchMaking.mu.Lock()
			app.ws.matchMaking.waiting = nil
			app.ws.matchMaking.mu.Unlock()
		}

		delete(app.ws.connUsers, oldConn)
		oldConn.Close()
	}

	app.ws.userConns[user.ID] = conn
	app.ws.connUsers[conn] = user.ID
	app.ws.userData[user.ID] = user
	app.ws.mu.Unlock()

	app.ws.matchPlayers(conn)
}

func (app *wsApp) matchPlayers(conn *websocket.Conn) {
	app.mu.Lock()
	currentUserID := app.connUsers[conn]
	currentUser := app.userData[currentUserID]
	app.mu.Unlock()

	if currentUserID == 0 || currentUser == nil {
		log.Println("Cannot identify user for this connection")
		return
	}

	app.matchMaking.mu.Lock()
	defer app.matchMaking.mu.Unlock()

	if app.matchMaking.waiting == nil {
		app.matchMaking.waiting = conn
		log.Printf("User %d waiting for opponent...\n", currentUserID)
		return
	}

	opponent := app.matchMaking.waiting
	app.matchMaking.waiting = nil

	app.mu.Lock()
	waitingUserID := app.connUsers[opponent]
	waitingUser := app.userData[waitingUserID]
	app.mu.Unlock()

	if waitingUserID == currentUserID || waitingUser == nil {
		log.Printf("User %d cannot match with themselves\n", currentUserID)
		msg := response{Type: "error", Message: "Cannot match with yourself. Please wait for another player."}
		msgJSON, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, msgJSON)
		return
	}

	question, err := app.app.store.Questions.GetRandomQuestion(context.Background())
	if err != nil {
		log.Println("Error fetching random question:", err)
		return
	}

	match := &Match{
		Player1:     conn,
		Player2:     opponent,
		Question:    *question,
		IsCompleted: false,
	}

	app.mu.Lock()
	app.matches[conn] = match
	app.matches[opponent] = match
	app.scores[conn] = 0
	app.scores[opponent] = 0
	app.mu.Unlock()

	log.Printf("Matched users %d and %d\n", currentUserID, waitingUserID)

	msgToCurrentUser := response{
		Type:    "question",
		Message: question,
		Opponent: &Opponent{
			Username: waitingUser.Username,
			Points:   waitingUser.Points,
		},
	}
	msgToOpponent := response{
		Type:    "question",
		Message: question,
		Opponent: &Opponent{
			Username: currentUser.Username,
			Points:   currentUser.Points,
		},
	}

	msgJSON1, _ := json.Marshal(msgToCurrentUser)
	msgJSON2, _ := json.Marshal(msgToOpponent)

	conn.WriteMessage(websocket.TextMessage, msgJSON1)
	opponent.WriteMessage(websocket.TextMessage, msgJSON2)

	go app.handleMessages(conn)
	go app.handleMessages(opponent)
}

func (app *wsApp) handleMessages(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var data payload
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		if data.Type != "answer" {
			continue
		}

		app.mu.Lock()
		match := app.matches[conn]
		app.mu.Unlock()

		if match == nil {
			log.Println("No active match found.")
			continue
		}

		match.mu.Lock()
		if match.IsCompleted {
			match.mu.Unlock()
			log.Println("Challenge already over")
			continue
		}
		match.mu.Unlock()

		stdin := match.Question.ExampleInput
		expectedOutput := strings.TrimSpace(match.Question.ExampleOutput)

		result, err := sendToJudge(data.Answer, data.LangID, stdin)
		if err != nil {
			log.Println("Judge0 error:", err)
			return
		}

		userOutput := strings.TrimSpace(result.Stdout)

		if normalizeOuput(userOutput) == normalizeOuput(expectedOutput) {
			match.mu.Lock()
			if !match.IsCompleted {
				match.IsCompleted = true

				app.mu.Lock()
				opponent := match.Player1
				if conn == match.Player1 {
					opponent = match.Player2
				}
				winnerID := app.connUsers[conn]
				loserID := app.connUsers[opponent]
				app.mu.Unlock()

				log.Println("Correct answer. Challenge over!")
				go app.app.updatePoints(winnerID, loserID, match.Question.ID)

				winMSG := response{Type: "feedback", Message: "Correct. You won!"}
				loseMSG := response{Type: "feedback", Message: "You lost!"}
				winJSON, _ := json.Marshal(winMSG)
				loseJSON, _ := json.Marshal(loseMSG)

				conn.WriteMessage(websocket.TextMessage, winJSON)
				opponent.WriteMessage(websocket.TextMessage, loseJSON)
			}
			match.mu.Unlock()
		} else {
			feedback := response{Type: "feedback", Message: "Incorrect!"}
			respJSON, _ := json.Marshal(feedback)
			conn.WriteMessage(websocket.TextMessage, respJSON)
		}
	}

	app.mu.Lock()
	match := app.matches[conn]
	userID := app.connUsers[conn]
	app.mu.Unlock()

	if match == nil {
		return
	}

	match.mu.Lock()
	if match.IsCompleted {
		match.mu.Unlock()
		return
	}
	match.IsCompleted = true
	match.mu.Unlock()

	var opponent *websocket.Conn
	app.mu.Lock()
	if match.Player1 == conn {
		opponent = match.Player2
	} else {
		opponent = match.Player1
	}
	opponentID := app.connUsers[opponent]
	app.mu.Unlock()

	winMSG := response{Type: "feedback", Message: "Your opponent disconnected. You won!"}
	winJSON, _ := json.Marshal(winMSG)

	opponent.WriteMessage(websocket.TextMessage, winJSON)

	go app.app.updatePoints(opponentID, userID, match.Question.ID)

	app.mu.Lock()
	delete(app.matches, conn)
	delete(app.matches, opponent)
	delete(app.scores, conn)
	delete(app.scores, opponent)
	app.mu.Unlock()
}

func sendToJudge(code string, langID int, stdin string) (submissionResponse, error) {
	reqBody, _ := json.Marshal(submissionRequest{
		LanguageID: langID,
		SourceCode: code,
		Stdin:      stdin,
	})

	req, err := http.NewRequest("POST", judge0APIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return submissionResponse{}, nil
	}

	req.Header.Set("X-RapidAPI-Key", env.GetString("RAPID_API_KEY", ""))
	req.Header.Set("X-RapidAPI-Host", env.GetString("RAPID_API_HOST", ""))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return submissionResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result submissionResponse
	json.Unmarshal(body, &result)

	return result, nil
}

func normalizeOuput(output string) string {
	output = strings.TrimSpace(output)
	output = strings.ReplaceAll(output, "\r\n", "\n")
	output = strings.Join(strings.Fields(output), " ")
	return output
}

func (app *application) updatePoints(winnerID, loserID, questionID int64) {
	ctx := context.Background()

	_, err := app.store.Users.IncrementPoints(ctx, winnerID, 10)
	if err != nil {
		log.Println("Error incrementing points:", err)
	}

	_, err = app.store.Users.DecrementPoints(ctx, loserID, 10)
	if err != nil {
		log.Println("Error decrementing points:", err)
	}

	matchResult := store.Match{
		Player1ID:  winnerID,
		Player2ID:  loserID,
		WinnerID:   winnerID,
		QuestionID: questionID,
	}

	err = app.store.Matches.Create(ctx, &matchResult)
	if err != nil {
		log.Println("Error storing match result:", err)
	}
}
