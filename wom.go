package wom

import (
	"bytes"
	"encoding/json"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"log"
	"net/http"
	"strings"
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

func checkGuess(app *pocketbase.PocketBase, r *models.Record) bool {
	content := strings.ToLower(r.GetString("content"))
	puzzle := r.Get("puzzle")
	var count string
	err := app.Dao().DB().Select("count(*)").From("answers").
		AndWhere(dbx.HashExp{"puzzle": puzzle}).
		AndWhere(dbx.HashExp{"content": content}).
		Row(&count)
	if err != nil {
		log.Printf("Error checking guess: %s", err)
		return false
	}
	return count == "1"
}
