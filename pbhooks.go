package wom

import (
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func ConfigurePocketBase(app *pocketbase.PocketBase) error {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	serveCmd := cmd.NewServeCommand(app, true)
	serveCmd.Flags().StringP("email", "e", viper.GetString("EMAIL"), "Sets the initial admin email")
	serveCmd.Flags().StringP("password", "p", viper.GetString("PASSWORD"), "Sets the initial admin password")
	serveCmd.Flags().StringP("webhook-url", "w", viper.GetString("WEBHOOK_URL"), "Webhook to send events to {'content': 'message'}")
	serveCmd.Flags().StringP("hcaptchaSecretKey", "", viper.GetString("HCAPTCHA_SECRET_KEY"), "Secret key to use for hCaptcha")
	serveCmd.Flags().StringP("hCaptchaSiteKey", "", viper.GetString("HCAPTCHA_SITE_KEY"), "Site key to use for hCaptcha")
	serveCmd.Flags().StringP("mailinglistSecretKey", "", viper.GetString("MAILINGLIST_SECRET_KEY"), "Mailing list secret key")
	err := serveCmd.MarkFlagRequired("hcaptchaSecretKey")
	if err != nil {
		return err
	}
	err = serveCmd.MarkFlagRequired("hCaptchaSiteKey")
	if err != nil {
		return err
	}
	err = serveCmd.MarkFlagRequired("mailinglistSecretKey")
	if err != nil {
		return err
	}
	app.RootCmd.AddCommand(serveCmd)
	app.RootCmd.AddCommand(cmd.NewAdminCommand(app))
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		return settings(app)
	})
	app.OnBeforeServe().Add(createAdminAccountHook(serveCmd))
	app.OnBeforeServe().Add(createWomRoutesHook(app))
	app.OnRecordBeforeUpdateRequest("adventures").Add(createPreserveFilenameUpdateHook)
	app.OnRecordBeforeCreateRequest("adventures").Add(createPreserveFilenameCreateHook)
	app.OnRecordBeforeCreateRequest("guesses").Add(createBeforeGuessCreatedHook(app))
	app.OnRecordAfterCreateRequest("guesses").Add(createGuessCreatedHook(app))
	return nil
}

func settings(app *pocketbase.PocketBase) error {
	key, _ := app.RootCmd.Flags().GetString("hcaptchaSecretKey")
	form := forms.NewSettingsUpsert(app)
	form.Meta.AppName = "Puzzad"
	form.Meta.AppUrl = "http://localhost.puzzad.com:8090"
	form.Meta.HideControls = true
	form.Logs.MaxDays = 90
	form.Smtp.Host = key

	return form.Submit()
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

func createGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		webhookURL, _ := app.RootCmd.Flags().GetString("webhook-url")
		game, err := app.Dao().FindRecordById("games", e.Record.Get("game").(string))
		if err != nil {
			return err
		}
		puzzle, err := app.Dao().FindRecordById("puzzles", e.Record.Get("puzzle").(string))
		if err != nil {
			return err
		}
		if e.Record.Get("correct").(bool) {
			sendWebhook(webhookURL, fmt.Sprintf(":tada: %s/%s: %s", game.Get("username"), puzzle.Get("title"), e.Record.Get("content")))
			if puzzle.Get("next") == "" {
				game.Set("puzzle", nil)
				game.Set("status", "EXPIRED")
				game.Set("end", time.Now())
			} else {
				game.Set("puzzle", puzzle.Get("next"))
			}
			return app.Dao().SaveRecord(game)
		}
		sendWebhook(webhookURL, fmt.Sprintf(":x: %s/%s: %s", game.Get("username"), puzzle.Get("title"), e.Record.Get("content")))
		return nil
	}
}

func createPreserveFilenameCreateHook(e *core.RecordCreateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
}

func createPreserveFilenameUpdateHook(e *core.RecordUpdateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
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
