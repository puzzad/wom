package wom

import (
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"time"
)

func ConfigurePocketBase(app *pocketbase.PocketBase, db *daos.Dao, mailClient mailer.Mailer, contactEmail, siteURL, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret, webhookURL string) {
	app.OnBeforeServe().Add(createWomRoutesHook(app, db, mailClient, webhookURL, contactEmail, siteURL, senderName, senderAddress, hcaptchaSecretKey, hcaptchaSiteKey, mailingListSecret))
	app.OnRecordBeforeUpdateRequest("adventures").Add(createPreserveFilenameUpdateHook)
	app.OnRecordBeforeCreateRequest("adventures").Add(createPreserveFilenameCreateHook)
	app.OnRecordBeforeCreateRequest("guesses").Add(createBeforeGuessCreatedHook(app))
	app.OnRecordAfterCreateRequest("guesses").Add(createGuessCreatedHook(app, webhookURL))
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

func createBeforeGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		e.Record.Set("correct", checkGuess(app, e.Record))
		return nil
	}
}

func createGuessCreatedHook(app *pocketbase.PocketBase, webhookURL string) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		game, err := app.Dao().FindRecordById("games", e.Record.Get("game").(string))
		if err != nil {
			return err
		}
		puzzle, err := app.Dao().FindRecordById("puzzles", e.Record.Get("puzzle").(string))
		if err != nil {
			return err
		}
		go func() {
			if e.Record.Get("correct").(bool) && puzzle.Get("next") == "" {
				sendWebhook(webhookURL, fmt.Sprintf(":tada: %s/%s: %s", game.Get("username"), puzzle.Get("title"), e.Record.Get("content")))
				sendWebhook(webhookURL, fmt.Sprintf(":checkered_flag:  %s finished", game.Get("username")))
			} else {
				sendWebhook(webhookURL, fmt.Sprintf(":x: %s/%s: %s", game.Get("username"), puzzle.Get("title"), e.Record.Get("content")))
			}
		}()
		if e.Record.Get("correct").(bool) {
			if puzzle.Get("next") == "" {
				game.Set("puzzle", nil)
				game.Set("status", "EXPIRED")
				game.Set("end", time.Now())
			} else {
				game.Set("puzzle", puzzle.Get("next"))
			}
			return app.Dao().SaveRecord(game)
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
