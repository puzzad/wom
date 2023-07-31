package main

import (
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/puzzad/wom"
	_ "github.com/puzzad/wom/migrations"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	dev   = flag.Bool("production", true, "Is this a dev instance")
	debug = flag.Bool("debug", false, "Enable debugging")

	adminEmail    = flag.String("email", "", "Sets the initial admin email")
	adminPassword = flag.String("password", "", "Sets the initial admin password")

	webhookURL = flag.String("webhook-url", "", "Webhook to send events to {'content': 'message'}")

	hcaptchaSecretKey    = flag.String("hcaptcha-secret-key", "", "Secret key to use for hCaptcha")
	hcaptchaSiteKey      = flag.String("hcaptcha-site-key", "", "Site key to use for hCaptcha")
	mailinglistSecretKey = flag.String("mailinglist-secret", "", "Mailing list secret key")

	smtpHost          = flag.String("smtp-host", "", "SMTP Host to send email via")
	smtpPort          = flag.Int("smtp-port", 25, "SMTP Port")
	smtpUser          = flag.String("smtp-user", "", "SMTP Username")
	smtpPass          = flag.String("smtp-pass", "", "SMTP password")
	smtpSenderAddress = flag.String("smtp-sender-email", "", "SMTP Sender address")
	smtpSenderName    = flag.String("smtp-sender-name", "", "SMTP Sender Name")

	siteURL      = flag.String("site-url", "", "Public facing site URL")
	siteName     = flag.String("site-name", "", "Public facing site name")
	backups      = flag.Bool("backups", false, "If enabled, backups will be performed every day at midnight, the last 7 will be kept")
	contactEmail = flag.String("contact-email", "", "Email address to send contact form emails to")

	createCollections = flag.Bool("create-migration", false, "Creates new migration file with snapshot of the local collections configuration.  Will write to /migrations")
	migrationSync     = flag.Bool("migration-sync", false, "Ensures that the _migrations history table doesn't have references to deleted migration files")
	autoMigrate       = flag.Bool("auto-migrate", false, "Automatically create migrations for actions taking in the admin UI.  Will write these to /migrations")
	dataDirectory     = flag.String("data-dir", "/data", "Directory to store database and backups")

	googleAuth    = flag.Bool("google-auth", false, "Enable Google oauth")
	googleID      = flag.String("google-id", "", "Google oauth client ID")
	googleSecret  = flag.String("google-secret", "", "Google oauth client secret")
	discordAuth   = flag.Bool("discord-auth", false, "Enable Discord oauth")
	discordID     = flag.String("discord-id", "", "Discord oauth client ID")
	discordSecret = flag.String("discord-secret", "", "Discord oauth client secret")
	twitchAuth    = flag.Bool("twitch-auth", false, "Enable Twitch oauth")
	twitchID      = flag.String("twitch-id", "", "Twitch oauth client ID")
	twitchSecret  = flag.String("twitch-secret", "", "Twitch oauth client secret")

	required = []string{
		"hcaptcha-site-key",
		"hcaptcha-secret-key",
		"mailinglist-secret",
		"site-name",
		"site-url",
		"smtp-sender-name",
		"smtp-sender-email",
	}
)

func main() {
	updateRequiredFlagsUsage()
	envflag.Parse()
	checkRequiredFlags()

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: *dataDirectory,
		DefaultDebug:   *debug,
	})

	if err := app.Bootstrap(); err != nil {
		log.Fatalf("Failed to bootstrap: %v", err)
	}

	if err := runMigrationCommand(app, "up"); err != nil {
		log.Fatalf("Unable to migrate: %s", err)
	}

	if err := UpdateSettings(app); err != nil {
		log.Fatalf("Failed to update settings: %v", err)
	}

	if err := UpdateAdmin(app, *adminEmail, *adminPassword); err != nil {
		log.Fatalf("Failed to set admin account: %v", err)
	}

	if *createCollections {
		err := runMigrationCommand(app, "collections")
		if err != nil {
			log.Fatalf("Failed to create collections migration: %v", err)
		}
		return
	}

	if *migrationSync {
		err := runMigrationCommand(app, "history-sync")
		if err != nil {
			log.Fatalf("Failed to sync migration history: %v", err)
		}
		return
	}

	wom.ConfigurePocketBase(
		app,
		app.Dao(),
		app.NewMailClient(),
		*contactEmail,
		*siteURL,
		app.Settings().Meta.SenderName,
		app.Settings().Meta.SenderAddress,
		*hcaptchaSecretKey,
		*hcaptchaSiteKey,
		*mailinglistSecretKey,
		*webhookURL,
		*dev,
	)

	serveCmd := cmd.NewServeCommand(app, false)
	serveCmd.SetArgs([]string{"--http=0.0.0.0:8090"})
	log.Printf("Starting wom: http://0.0.0.0:8090/_/")
	if err := serveCmd.Execute(); err != nil {
		log.Fatalf("Error whilst serving: %v", err)
	}
}

func updateRequiredFlagsUsage() {
	flag.VisitAll(func(f *flag.Flag) {
		if contains(required, f.Name) {
			f.Usage = fmt.Sprintf("%s [REQUIRED]", f.Usage)
		}
	})
}

func checkRequiredFlags() {
	seen := make(map[string]bool)
	flag.VisitAll(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	for _, r := range required {
		if !seen[r] {
			log.Fatalf("Missing required flag: %s", r)
		}
	}
}

func runMigrationCommand(app core.App, subcommand string) error {
	blankCommand := &cobra.Command{}
	migratecmd.MustRegister(app, blankCommand, migratecmd.Config{Automigrate: *autoMigrate})
	blankCommand.SetArgs([]string{"migrate", subcommand})
	return blankCommand.Execute()
}

func UpdateAdmin(app *pocketbase.PocketBase, email, password string) error {
	if len(email) == 0 && len(password) == 0 {
		return nil
	}
	if is.EmailFormat.Validate(email) != nil || len(password) <= 5 {
		return fmt.Errorf("invalid admin credentials\n")
	}
	admin, err := app.Dao().FindAdminByEmail(email)
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
	err = app.Dao().SaveAdmin(admin)
	if err != nil {
		fmt.Printf("Error saving admin: %v\n", err)
	}
	return err
}

func UpdateSettings(app *pocketbase.PocketBase) error {
	form := forms.NewSettingsUpsert(app)
	if strings.HasSuffix(*siteName, "/") {
		form.Meta.AppName = *siteName
	} else {
		form.Meta.AppName = fmt.Sprintf("%s/", *siteName)
	}
	form.Meta.AppUrl = *siteURL
	form.Meta.HideControls = true
	form.Logs.MaxDays = 90
	form.Smtp.Enabled = *smtpHost != ""
	if *smtpHost != "" {
		form.Smtp.Host = *smtpHost
		form.Smtp.Port = *smtpPort
		form.Smtp.Username = *smtpUser
		form.Smtp.Password = *smtpPass
	}
	form.Meta.SenderName = *smtpSenderName
	form.Meta.SenderAddress = *smtpSenderAddress
	if *backups {
		form.Backups.Cron = "0 0 * * *"
		form.Backups.CronMaxKeep = 7
	}

	form.GoogleAuth.Enabled = *googleAuth
	form.GoogleAuth.ClientId = *googleID
	form.GoogleAuth.ClientSecret = *googleSecret

	form.DiscordAuth.Enabled = *discordAuth
	form.DiscordAuth.ClientId = *discordID
	form.DiscordAuth.ClientSecret = *discordSecret

	form.TwitchAuth.Enabled = *twitchAuth
	form.TwitchAuth.ClientId = *twitchID
	form.TwitchAuth.ClientSecret = *twitchSecret

	//TODO: Add frontend endpoints for these so they fit in with the general site theme rather than point at pocketbase
	template, err := wom.GetTemplates("changeemail", wom.TemplateData{})
	if err != nil {
		return err
	}
	form.Meta.ConfirmEmailChangeTemplate.Subject = "{APP_NAME}: Confirm email change"
	form.Meta.ConfirmEmailChangeTemplate.ActionUrl = "{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}"
	form.Meta.ConfirmEmailChangeTemplate.Body = template
	template, err = wom.GetTemplates("resetpassword", wom.TemplateData{})
	if err != nil {
		return err
	}
	form.Meta.ResetPasswordTemplate.Subject = "{APP_NAME}: Password reset"
	form.Meta.ConfirmEmailChangeTemplate.ActionUrl = "{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}"
	form.Meta.ResetPasswordTemplate.Body = template
	template, err = wom.GetTemplates("verification", wom.TemplateData{})
	if err != nil {
		return err
	}
	form.Meta.VerificationTemplate.Subject = "{APP_NAME}: Email Verification"
	form.Meta.ConfirmEmailChangeTemplate.ActionUrl = "{APP_URL}/_/#/auth/confirm-verification/{TOKEN}"
	form.Meta.VerificationTemplate.Body = template
	return form.Submit()
}

func contains[T comparable](s []T, e T) bool {
	//TODO: Remove in go 1.21
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
