package wom

import (
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"net/mail"
	"strings"
)

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
