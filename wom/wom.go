package wom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/csmith/aca"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func sendWebhook(webHookURL, message string) {
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
}

func checkGuess(app *pocketbase.PocketBase, r *models.Record) bool {
	content := r.Get("content")
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

func startGame(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(c echo.Context) error {
		code := c.PathParam("code")
		if len(code) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Game"})
		}
		q := app.Dao().DB().NewQuery("UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures.firstpuzzle FROM adventures WHERE adventures.id = games.adventure), start = datetime('now') WHERE username = {:username} AND status = 'PAID' AND (puzzle='' OR puzzle IS NULL);")
		q = q.Bind(dbx.Params{"username": code})
		result, err := q.Execute()
		if err != nil {
			fmt.Printf("%v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 1"})
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 2"})
		}
		if rows != 1 {
			fmt.Printf("Rows: %d\n", rows)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 3"})
		}
		return nil
	}
}

func startAdventure(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")
		if len(id) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Adventure"})
		}
		adventure, err := app.Dao().FindRecordById("adventures", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Adventure not found"})
		}
		collection, err := app.Dao().FindCollectionByNameOrId("games")
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Adventure not found"})
		}
		user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if user == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
		}
		if !user.Verified() {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email must be verified"})
		}
		acaGen, err := aca.NewGenerator(".", rand.NewSource(time.Now().UnixNano()))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to generate ACA"})
		}
		code := acaGen.Generate()
		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		record.RefreshId()
		err = form.LoadData(map[string]any{
			"status":          "PAID",
			"user":            user.Id,
			"adventure":       adventure.Id,
			"username":        code,
			"password":        "puzzad",
			"passwordConfirm": "puzzad",
		})
		if err = form.Submit(); err != nil {
			fmt.Printf("Unable to add Adventure: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		currentGames := user.GetStringSlice("games")
		currentGames = append(currentGames, record.Id)
		user.Set("games", currentGames)
		err = app.Dao().SaveRecord(user)
		if err != nil {
			fmt.Printf("Unable to add Adventure: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add game to user"})
		}
		return c.JSON(http.StatusOK, map[string]string{"code": code})
	}
}
