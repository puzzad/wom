package wom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/csmith/aca"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/spf13/cobra"
	"math/rand"
	"net/http"
	"time"
)

func ConfigurePocketBase(app *pocketbase.PocketBase) {
	app.OnBeforeBootstrap().Add(func(e *core.BootstrapEvent) error {
		var serveCmd *cobra.Command
		for i := range app.RootCmd.Commands() {
			if app.RootCmd.Commands()[i].Name() == "serve" {
				serveCmd = app.RootCmd.Commands()[i]
			}
		}
		if serveCmd == nil {
			return fmt.Errorf("unable to find serve command")
		}
		serveCmd.Flags().StringP("email", "e", "", "Sets the initial admin email")
		serveCmd.Flags().StringP("password", "p", "", "Sets the initial admin password")
		serveCmd.Flags().StringP("webhook-url", "w", "", "Webhook to send events to {'content': 'message'}")
		app.RootCmd.AddCommand(NewImportCmd())
		app.OnBeforeServe().Add(createAdminAccountHook(serveCmd))
		app.OnBeforeServe().Add(createWomRoutesHook(app))
		app.OnRecordBeforeUpdateRequest("adventures").Add(createPreserveFilenameUpdateHook)
		app.OnRecordBeforeCreateRequest("adventures").Add(createPreserveFilenameCreateHook)
		app.OnRecordBeforeCreateRequest("guesses").Add(createBeforeGuessCreatedHook(app))
		app.OnRecordAfterCreateRequest("guesses").Add(createGuessCreatedHook(app))
		return nil
	})
}

func createBeforeGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		e.Record.Set("correct", checkGuess(app, e.Record))
		return nil
	}
}

func createGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		var code, title string
		err := app.Dao().DB().Select("games.username as username", "puzzles.title as title").
			From("guesses").
			InnerJoin("games", dbx.NewExp("games.id=guesses.game")).
			InnerJoin("puzzles", dbx.NewExp("puzzles.id=guesses.puzzle")).
			Where(dbx.HashExp{"guesses.id": e.Record.Id}).
			Row(&code, &title)
		if err == nil {
			webhookURL, _ := app.RootCmd.Flags().GetString("webhook-url")
			if e.Record.Get("correct").(bool) {
				sendWebhook(webhookURL, fmt.Sprintf(":tada: %s/%s: %s", code, title, e.Record.Get("content")))
			} else {
				sendWebhook(webhookURL, fmt.Sprintf(":x: %s/%s: %s", code, title, e.Record.Get("content")))
			}
		}
		return nil
	}
}

func createPreserveFilenameCreateHook(e *core.RecordCreateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
}

func createPreserveFilenameUpdateHook(e *core.RecordUpdateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
}

func createWomRoutesHook(app *pocketbase.PocketBase) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		_, err := e.Router.AddRoute(echo.Route{
			Name:    "start adventure",
			Path:    "/adventure/:id/start",
			Method:  http.MethodPost,
			Handler: startAdventure(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Name:    "start game",
			Path:    "/games/:code/start",
			Method:  http.MethodPost,
			Handler: startGame(app),
		})
		if err != nil {
			return err
		}
		_, err = e.Router.AddRoute(echo.Route{
			Name:   "send contact form",
			Path:   "/mail/contact",
			Method: http.MethodPost,
			//Handler: wom.SendContactForm,
		})
		//e.Router.Add(http.MethodGet, "/mail/subscribe", wom.SubscribeToMailingList)
		//e.Router.Add(http.MethodGet, "/mail/confirm", wom.ConfirmMailingListSubscription)
		//e.Router.Add(http.MethodGet, "/mail/unsubscribe", wom.UnsubscribeFromMailingList)
		//e.Router.Add(http.MethodGet, "/mail/contact", wom.SendContactForm)
		return nil
	}
}

func createAdminAccountHook(serveCmd *cobra.Command) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		email, _ := serveCmd.Flags().GetString("email")
		password, _ := serveCmd.Flags().GetString("password")
		if len(email) == 0 && len(password) == 0 {
			return nil
		}
		if is.EmailFormat.Validate(email) != nil || len(password) <= 5 {
			return fmt.Errorf("invalid admin credentials\n")
		}
		admin, err := e.App.Dao().FindAdminByEmail(email)
		if err != nil {
			fmt.Printf("Creating admin account: %s\n", email)
			admin = &models.Admin{
				Email: email,
			}
		}
		err = admin.SetPassword(password)
		if err != nil {
			fmt.Printf("Error setting admin password: %v\n", err)
			return err
		}
		err = e.App.Dao().SaveAdmin(admin)
		if err != nil {
			fmt.Printf("Error saving admin: %v\n", err)
		}
		return err
	}
}

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
	game := r.Get("game")
	var count string
	err := app.Dao().DB().Select("count(*)").From("answers").
		AndWhere(dbx.HashExp{"puzzle": puzzle}).
		AndWhere(dbx.HashExp{"game": game}).
		AndWhere(dbx.HashExp{"answer": content}).
		Row(&count)
	if err != nil {
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
		q := app.Dao().DB().NewQuery("UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures.firstpuzzle FROM adventures WHERE adventures.id = games.adventure), start = datetime('now') WHERE code = {:code} AND status = 'PAID' AND (puzzle='' OR puzzle IS NULL);")
		q = q.Bind(dbx.Params{"code": code})
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
		acaGen, err := aca.NewGenerator("-", rand.NewSource(time.Now().UnixNano()))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to generate ACA"})
		}
		code := acaGen.Generate()
		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		err = form.LoadData(map[string]any{
			"status":    "PAID",
			"user":      user.Id,
			"adventure": adventure.Id,
			"code":      code,
		})
		if err = form.Submit(); err != nil {
			fmt.Printf("%v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		if err = form.Submit(); err != nil {
			fmt.Printf("%v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		return c.JSON(http.StatusOK, map[string]string{"code": code})
	}
}

func preserveOriginalFilenames(uploadedFiles map[string][]*filesystem.File, record *models.Record) error {
	oldNames := map[string]map[string]string{}
	for field, files := range uploadedFiles {
		if len(files) == 0 {
			continue
		}
		oldNames[field] = make(map[string]string, len(files))
		for _, f := range files {
			oldNames[field][f.Name] = f.OriginalName
			f.Name = f.OriginalName
		}
	}
	for field, filenames := range oldNames {
		files := record.GetStringSlice(field)

		for i, old := range files {
			if newName, ok := filenames[old]; ok {
				files[i] = newName
			}
		}
		record.Set(field, files)
	}
	return nil
}
