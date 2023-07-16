package wom

import (
	"bytes"
	"encoding/json"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"log"
	"net/http"
)

func sendWebhook(webHookURL, message string) {
	go func() {
		type webhook struct {
			Content string `json:"content"`
		}
		if len(webHookURL) == 0 {
			return
		}
		data := &webhook{
			Content: message,
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return
		}
		_, err = http.Post(webHookURL, "application/json", bytes.NewReader(dataBytes))
		if err != nil {
			return
		}
	}()
}

func checkGuess(db *daos.Dao, puzzleID, guess string) bool {
	records, err := db.FindRecordsByExpr("answers", dbx.HashExp{"puzzle": puzzleID, "content": guess})
	if err != nil {
		log.Printf("Error checking guess: %s", err)
		return false
	}
	return len(records) == 1
}
