package main

import (
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/puzzad/wom"
	_ "github.com/puzzad/wom/migrations"
	"log"
)

var (
	adminEmail           = flag.String("email", "", "Sets the initial admin email")
	adminPassword        = flag.String("password", "", "Sets the initial admin password")
	webhookURL           = flag.String("webhook-url", "", "Webhook to send events to {'content': 'message'}")
	hcaptchatSecretKey   = flag.String("hcaptcha-secret-key", "", "Secret key to use for hCaptcha")
	hcaptchaSiteKey      = flag.String("hcaptcha-site-key", "", "Site key to use for hCaptcha")
	mailinglistSecretKey = flag.String("mailinglist-secret", "", "Mailing list secret key")
	smtpHost             = flag.String("smtp-host", "", "SMTP Host to send email via")
	smtpPort             = flag.Int("smtp-port", 25, "SMTP Port")
	smtpUser             = flag.String("smtp-user", "", "SMTP Username")
	smtpPass             = flag.String("smtp-pass", "", "SMTP password")
	smtpSenderAddress    = flag.String("smtp-sender-email", "", "SMTP Sender address")
	smtpSenderName       = flag.String("smtp-sender-name", "", "SMTP Sender Name")
	siteURL              = flag.String("site-url", "", "Public facing site URL")
	siteName             = flag.String("site-name", "", "Public facing site name")
	backups              = flag.Bool("backups", false, "If enabled, backups will be performed every day at midnight, the last 7 will be kept")
	contactEmail         = flag.String("contact-email", "", "Email address to send contact form emails to")
)

func main() {
	envflag.Parse()
	if *hcaptchaSiteKey == "" || *hcaptchatSecretKey == "" || *mailinglistSecretKey == "" {
		log.Fatal("Missing required flags")
	}
	app := pocketbase.NewWithConfig(&pocketbase.Config{
		DefaultDataDir: "./data",
	})
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})
	if err := app.Bootstrap(); err != nil {
		log.Fatal(err)
	}
	if err := UpdateSettings(app); err != nil {
		log.Fatal(err)
	}
	if err := UpdateAdmin(app, *adminEmail, *adminPassword); err != nil {
		log.Fatal(err)
	}
	wom.ConfigurePocketBase(app, app.Dao(), app.NewMailClient(), *contactEmail, *siteURL, app.Settings().Meta.SenderName,
		app.Settings().Meta.SenderAddress, *hcaptchatSecretKey, *hcaptchaSiteKey, *mailinglistSecretKey, *webhookURL)
	serveCmd := cmd.NewServeCommand(app, false)
	serveCmd.SetArgs([]string{"--http=0.0.0.0:8090"})
	log.Printf("Starting wom: http://0.0.0.0:8090/_/")
	if err := serveCmd.Execute(); err != nil {
		log.Fatal(err)
	}
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
	form.Meta.AppName = *siteName
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
	return form.Submit()
}
