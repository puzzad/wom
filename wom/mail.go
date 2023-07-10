package wom

import (
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"net/http"
	"net/mail"
	"strings"
)

func sendContactForm(app *pocketbase.PocketBase) echo.HandlerFunc {
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
			secretKey, _ := app.RootCmd.Flags().GetString("hcaptchaSecretKey")
			siteKey, _ := app.RootCmd.Flags().GetString("hcaptchaSiteKey")
			if err := checkCaptcha(siteKey, secretKey, data.Token); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
			if !validEmail(data.Email) {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
		}
		if err := SendContactFormMail(app, data.Email, data.Name, data.Message); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func SendContactFormMail(app *pocketbase.PocketBase, email string, name string, content string) error {
	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Headers: map[string]string{
			"Reply-To": email,
		},
		Subject: "Contact Form",
		Text:    fmt.Sprintf("Name: %s\nMessage\n%s", name, content),
	}
	return app.NewMailClient().Send(message)
}

func sendSubscriptionConfirmedMail(app *pocketbase.PocketBase, email string) error {
	secretKey, _ := app.RootCmd.Flags().GetString("mailinglistSecretKey")
	token, err := createUnsubscribeJwt(secretKey, email)
	if err != nil {
		return err
	}
	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Confirmed",
		Text:    fmt.Sprintf("/mail/unsubscribe/%s", token),
	}
	return app.NewMailClient().Send(message)
}

func sendSubscriptionUnsubscribedMail(app *pocketbase.PocketBase, email string) error {
	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Unsubscribed",
		Text:    fmt.Sprintf("Sorry to see you go"),
	}
	return app.NewMailClient().Send(message)
}

func sendSubscriptionOptInMail(app *pocketbase.PocketBase, email string) error {
	secretKey, _ := app.RootCmd.Flags().GetString("mailinglistSecretKey")
	token, err := createSubscriptionJwt(secretKey, email)
	if err != nil {
		return err
	}
	message := &mailer.Message{
		From: mail.Address{
			Name:    app.Settings().Meta.SenderName,
			Address: app.Settings().Meta.SenderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Opt-In",
		Text:    fmt.Sprintf("/mail/confirm/%s", token),
	}
	return app.NewMailClient().Send(message)
}

func addEmailToMailingList(app *pocketbase.PocketBase, email string) error {
	mailinglist, err := app.Dao().FindCollectionByNameOrId("mailinglist")
	if err != nil {
		return err
	}
	record := models.NewRecord(mailinglist)
	record.RefreshId()
	record.Set("email", email)
	return app.Dao().SaveRecord(record)
}

func removeEmailToMailingList(app *pocketbase.PocketBase, email string) error {
	record, err := app.Dao().FindFirstRecordByData("mailinglist", "email", email)
	if err != nil {
		return err
	}
	return app.Dao().DeleteRecord(record)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && !strings.Contains(email, " ")
}
