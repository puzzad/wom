package wom

import (
	"archive/zip"
	"embed"
	"fmt"
	"github.com/csmith/aca"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//go:embed all:site
var embedded embed.FS

var siteFS = echo.MustSubFS(embedded, "site")

func createWomRoutesHook(app core.App, fso fileSystemOpener, db *daos.Dao, mailClient mailer.Mailer, webhookURL, contactEmail, siteURL, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret string, dev bool) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		_ = e.Router.POST("/wom/signup", handleSignup(app, db, webhookURL))
		_ = e.Router.POST("/wom/startadventure", handleStartAdventure(db, webhookURL))
		_ = e.Router.POST("/wom/startgame", handleStartGame(db, webhookURL))
		_ = e.Router.POST("/wom/importzip", handleAdventureImport(db, fso, dev), apis.RequireAdminAuth())
		_ = e.Router.POST("/wom/requesthint", handleHintRequest(db, webhookURL))
		_ = e.Router.GET("/wom/gethints", getHints(db))
		_ = e.Router.POST("/wom/contact", handleContactForm(mailClient, contactEmail, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey))
		_ = e.Router.POST("/wom/subscribe", handleSubscribe(db, mailClient, siteURL, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret))
		_ = e.Router.GET("/wom/confirm/:token", handleConfirm(db, mailClient, siteURL, senderName, senderAddress, mailingListSecret))
		_ = e.Router.GET("/wom/unsubscribe/:token", handleUnsubscribe(db, mailClient, siteURL, senderName, senderAddress, mailingListSecret))
		e.Router.GET("/*", apis.StaticDirectoryHandler(siteFS, true))
		return nil
	}
}

func handleSignup(app core.App, db *daos.Dao, webhookURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		}
		if !validEmail(data.Email) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email"})
		}
		if !validPassword(data.Password) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid password"})
		}
		if _, err := db.FindAuthRecordByEmail("users", data.Email); err == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email"})
		}
		users, err := db.FindCollectionByNameOrId("users")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		record := models.NewRecord(users)
		err = record.SetEmail(data.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		err = record.SetPassword(data.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		err = record.SetUsername(data.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		if err = db.SaveRecord(record); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		if err = mails.SendRecordVerification(app, record); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to create user"})
		}
		sendWebhook(webhookURL, fmt.Sprintf("New user signup: `%s`", record.Id))
		return c.JSON(http.StatusOK, "User created")
	}
}

func getHints(db *daos.Dao) echo.HandlerFunc {
	return func(c echo.Context) error {
		type res struct {
			Id      string `json:"id"`
			Title   string `json:"title"`
			Message string `json:"message"`
		}
		game, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if game == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Game not found"})
		}
		puzzle := game.GetString("puzzle")
		if puzzle == "" {
			return c.NoContent(http.StatusNoContent)
		}
		var data []res
		err := db.DB().
			Select("hints.id as id", "hints.title as title", "IIF(usedhints.id ISNULL, '', hints.message) as message").
			From("hints").
			LeftJoin("usedhints", dbx.NewExp("hints.id=usedhints.hint AND usedhints.game={:gameID}", dbx.Params{"gameID": game.Id})).
			AndWhere(dbx.HashExp{"hints.puzzle": puzzle}).
			OrderBy("hints.order").
			All(&data)
		if err != nil {
			log.Printf("Error getting hints: %s", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error getting hints"})
		}
		return c.JSON(http.StatusOK, data)
	}
}

func handleStartAdventure(db *daos.Dao, webhookURL string) func(echo.Context) error {
	return func(c echo.Context) error {
		type req struct {
			Adventure string `json:"adventure"`
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		adventure, err := db.FindRecordById("adventures", data.Adventure)
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
			log.Printf("Failed to create ACA generator: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		code := acaGen.Generate()
		record := models.NewRecord(collection)
		record.RefreshId()
		record.Set("status", "PAID")
		record.Set("purchaser", user.Id)
		record.Set("adventure", adventure.Id)
		record.Set("username", code)
		if err := record.SetPassword("puzzad"); err != nil {
			log.Printf("Failed to set password on new game: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		if err := db.SaveRecord(record); err != nil {
			log.Printf("Failed to save adventure: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		sendWebhook(webhookURL, fmt.Sprintf("New game created: `%s` (adventure: %s)", code, adventure.Get("name")))
		currentGames := user.GetStringSlice("games")
		currentGames = append(currentGames, record.Id)
		user.Set("games", currentGames)
		err = db.SaveRecord(user)
		if err != nil {
			log.Printf("Unable to add game to user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add game to user"})
		}
		return c.JSON(http.StatusOK, map[string]string{"code": code})
	}
}

func handleStartGame(db *daos.Dao, webhookURL string) func(echo.Context) error {
	return func(c echo.Context) error {
		game, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if game == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Game not found"})
		}
		if game.GetString("puzzle") != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Game already started"})
		}
		adventure, err := db.FindFirstRecordByData("adventures", "id", game.Get("adventure"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to find first puzzle"})
		}
		game.Set("status", "ACTIVE")
		game.Set("puzzle", adventure.Get("firstpuzzle"))
		game.Set("start", time.Now())
		err = db.SaveRecord(game)
		if err != nil {
			log.Printf("Unable to start game: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game"})
		}
		sendWebhook(webhookURL, fmt.Sprintf(":rocket: `%s` started", game.Get("username")))
		return nil
	}
}

func handleContactForm(mailClient mailer.Mailer, contactEmail, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token   string `json:"token"`
			Name    string `json:"name"`
			Email   string `json:"email"`
			Message string `json:"message"`
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
		if err := SendContactFormMail(mailClient, contactEmail, senderName, senderAddress, data.Email, data.Name, data.Message); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleAdventureImport(db *daos.Dao, fso fileSystemOpener, dev bool) func(c echo.Context) error {
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
		adventures := getAdventures(zipReader, dev)
		updateAdventures(db, fso, adventures)
		return c.JSON(http.StatusOK, nil)
	}
}

func handleHintRequest(db *daos.Dao, webhookURL string) echo.HandlerFunc {
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
		type req struct {
			Hint string `json:"hint"`
		}
		var data = req{}
		if err = c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Hint ID not found"})
		}
		record := models.NewRecord(usedhints)
		record.RefreshId()
		record.Set("hint", data.Hint)
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
		go func() {
			hint, err := db.FindRecordById("hints", data.Hint)
			if err != nil {
				return
			}
			puzzle, err := db.FindRecordById("puzzles", game.GetString("puzzle"))
			sendWebhook(webhookURL, fmt.Sprintf(":bulb: `%s`/`%s`: hint requested: `%s`", game.Get("username"), puzzle.Get("title"), hint.Get("title")))
		}()
		return c.JSON(http.StatusOK, "")
	}
}

func handleSubscribe(db *daos.Dao, mailClient mailer.Mailer, siteURL, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		authInfo, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		type req struct {
			Email   string `json:"email"`
			Captcha string `json:"captcha"`
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email := ""
		if authInfo != nil && authInfo.Verified() {
			email = authInfo.Email()
		}
		if data.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if data.Email != email {
			if err := checkCaptcha(hcaptchaSiteKey, hcaptchaSecretKey, data.Captcha); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
			if err := sendSubscriptionOptInMail(mailClient, siteURL, senderName, senderAddress, mailingListSecret, data.Email); err != nil {
				fmt.Printf("Unable to send subscription email to '%s': %s\n", email, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error sending the opt in email"})
			}
			return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": true})
		}
		if err := addEmailToMailingList(db, data.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error adding you to the mailing list"})
		}
		if err := sendSubscriptionConfirmedMail(mailClient, siteURL, senderName, senderAddress, mailingListSecret, data.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error sending the confirmation email"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": false})
	}
}

func handleConfirm(db *daos.Dao, mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string `param:"token"`
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email, err := validateSubscriptionJWT(mailingListSecret, "subscribe", data.Token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err = addEmailToMailingList(db, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err = sendSubscriptionConfirmedMail(mailClient, siteURL, senderName, senderAddress, mailingListSecret, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleUnsubscribe(db *daos.Dao, mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string `param:"token"`
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
		if err := sendSubscriptionUnsubscribedMail(mailClient, siteURL, senderName, senderAddress, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
