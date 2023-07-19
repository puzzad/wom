package wom

import (
	"bytes"
	"encoding/json"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"log"
	"net/http"
	"unicode"
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

func validPassword(password string) bool {
	var length, upperCase, lowerCase, number, special bool
	if len(password) >= 8 {
		length = true
	}
	for i := range password {
		switch {
		case unicode.IsUpper(rune(password[i])):
			upperCase = true
		case unicode.IsLower(rune(password[i])):
			lowerCase = true
		case unicode.IsNumber(rune(password[i])):
			number = true
		case unicode.IsPunct(rune(password[i])) || unicode.IsSymbol(rune(password[i])):
			special = true
		}
	}
	return length && upperCase && lowerCase && number && special
}
