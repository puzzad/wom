package wom

import (
	"archive/zip"
	"fmt"
	"github.com/csmith/aca"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func createWomRoutesHook(app *pocketbase.PocketBase, db *daos.Dao, mailClient mailer.Mailer, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret string) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		_ = e.Router.POST("/adventure/:id/start", handleStartAdventure(db, app))
		_ = e.Router.POST("/games/:code/start", handleStartGame(db))
		_ = e.Router.POST("/mail/contact", handleContactForm(mailClient, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey))
		_ = e.Router.POST("/import/zip", handleAdventureImport(db, app), apis.RequireAdminAuth())
		_ = e.Router.POST("/hints/request", handleHintRequest(db))
		_ = e.Router.POST("/mail/subscribe", handleSubscribe(db, mailClient, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret))
		_ = e.Router.POST("/mail/confirm", handleConfirm(db, mailClient, senderName, senderAddress, mailingListSecret))
		_ = e.Router.POST("/mail/unsubscribe", handleUnsubscribe(db, mailClient, senderName, senderAddress, mailingListSecret))
		return nil
	}
}

func handleStartAdventure(db *daos.Dao, app core.App) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")
		if len(id) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Adventure"})
		}
		adventure, err := db.FindRecordById("adventures", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Adventure not found"})
		}
		collection, err := db.FindCollectionByNameOrId("games")
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

func handleStartGame(db *daos.Dao) func(echo.Context) error {
	return func(c echo.Context) error {
		code := c.PathParam("code")
		if len(code) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Game"})
		}
		q := db.DB().NewQuery("UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures.firstpuzzle FROM adventures WHERE adventures.id = games.adventure), start = datetime('now') WHERE username = {:username} AND status = 'PAID' AND (puzzle='' OR puzzle IS NULL);")
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

func handleContactForm(mailClient mailer.Mailer, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token   string
			Name    string
			Email   string
			Message string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if strings.TrimSpace(data.Name) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if strings.TrimSpace(data.Message) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		authInfo, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		email := ""
		if authInfo != nil && authInfo.Verified() {
			email = authInfo.Email()
		}
		if data.Email != email {
			if err := checkCaptcha(hcaptchaSiteKey, hcaptchaSecretKey, data.Token); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
			if !validEmail(data.Email) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
		}
		if err := SendContactFormMail(mailClient, senderName, senderAddress, data.Email, data.Name, data.Message); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleAdventureImport(db *daos.Dao, fso fileSystemOpener) func(c echo.Context) error {
	return func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		if len(form.File["adventures.zip"]) != 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Must be one file called adventures.zip"})
		}
		file := form.File["adventures.zip"][0]
		fileReader, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to open file"})
		}
		zipReader, err := zip.NewReader(fileReader, file.Size)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to open zip"})
		}
		adventures := getAdventures(zipReader, false)
		updateAdventures(db, fso, adventures)
		return c.JSON(http.StatusOK, nil)
	}
}

func handleHintRequest(db *daos.Dao) echo.HandlerFunc {
	return func(c echo.Context) error {
		game, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if game == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Game not found"})
		}
		usedhints, err := db.FindCollectionByNameOrId("usedhints")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "usedhints not found"})
		}
		guesses, err := db.FindCollectionByNameOrId("guesses")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "guesses not found"})
		}
		record := models.NewRecord(usedhints)
		record.RefreshId()
		hintID, err := io.ReadAll(c.Request().Body)
		defer func() {
			_ = c.Request().Body.Close()
		}()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "hintID not found"})
		}
		record.Set("hint", hintID)
		record.Set("game", game.Id)
		err = db.SaveRecord(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to request hint"})
		}
		record = models.NewRecord(guesses)
		record.RefreshId()
		record.Set("game", game.Id)
		record.Set("puzzle", game.Get("puzzle"))
		record.Set("content", "*hint")
		err = db.SaveRecord(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to request hint"})
		}
		return c.JSON(http.StatusOK, "")
	}
}

func handleSubscribe(db *daos.Dao, mailClient mailer.Mailer, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		authInfo, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		type req struct {
			Email   string
			Captcha string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email := ""
		if authInfo != nil && authInfo.Verified() {
			email = authInfo.Email()
		}
		if data.Email != email {
			if err := checkCaptcha(hcaptchaSiteKey, hcaptchaSecretKey, data.Captcha); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
			if err := sendSubscriptionOptInMail(mailClient, senderName, senderAddress, mailingListSecret, email); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
			}
			return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": true})
		}
		if err := addEmailToMailingList(db, data.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionConfirmedMail(mailClient, senderName, senderAddress, mailingListSecret, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": false})
	}
}

func handleConfirm(db *daos.Dao, mailClient mailer.Mailer, senderName, senderAddress, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email, err := validateSubscriptionJWT(mailingListSecret, "unsubscribe", data.Token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err := removeEmailToMailingList(db, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionConfirmedMail(mailClient, senderName, senderAddress, mailingListSecret, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleUnsubscribe(db *daos.Dao, mailClient mailer.Mailer, senderName, senderAddress, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email, err := validateSubscriptionJWT(mailingListSecret, "unsubscribe", data.Token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err := removeEmailToMailingList(db, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionUnsubscribedMail(mailClient, senderName, senderAddress, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
