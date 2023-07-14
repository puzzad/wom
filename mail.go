package wom

import (
	"bytes"
	"fmt"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	htemplate "html/template"
	"net/mail"
	"path/filepath"
	"strings"
	ttemplate "text/template"
)

func SendContactFormMail(mailClient mailer.Mailer, contactEmail, senderName, senderAddress, email string, name string, content string) error {
	text, html, err := getTemplates("subscribed", templateData{
		Name:    name,
		Email:   email,
		Message: content,
	})
	if err != nil {
		return err
	}
	message := &mailer.Message{
		From: mail.Address{
			Name:    senderName,
			Address: senderAddress,
		},
		To: []mail.Address{{
			Address: contactEmail,
		}},
		Headers: map[string]string{
			"Reply-To": email,
		},
		Subject: "Contact Form",
		Text:    text,
		HTML:    html,
	}
	return mailClient.Send(message)
}

func sendSubscriptionConfirmedMail(mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createUnsubscribeJwt(mailingListSecret, email)
	if err != nil {
		return err
	}
	text, html, err := getTemplates("subscribed", templateData{
		Link: fmt.Sprintf("%s/mail/unsubscribe/%s", siteURL, token),
	})
	if err != nil {
		return err
	}
	message := &mailer.Message{
		From: mail.Address{
			Name:    senderName,
			Address: senderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Confirmed",
		Text:    text,
		HTML:    html,
	}
	return mailClient.Send(message)
}

func sendSubscriptionUnsubscribedMail(mailClient mailer.Mailer, siteURL, senderName, senderAddress, email string) error {
	text, html, err := getTemplates("unsubscribed", templateData{
		Link: fmt.Sprintf("%s/mail/subscribe", siteURL),
	})
	if err != nil {
		return err
	}
	message := &mailer.Message{
		From: mail.Address{
			Name:    senderName,
			Address: senderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Unsubscribed",
		Text:    text,
		HTML:    html,
	}
	return mailClient.Send(message)
}

func sendSubscriptionOptInMail(mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createSubscriptionJwt(mailingListSecret, email)
	if err != nil {
		fmt.Printf("Error creating subscription JWT: %s", err)
		return err
	}
	text, html, err := getTemplates("optin", templateData{
		Link: fmt.Sprintf("%s/mail/confirm/%s", token, siteURL),
	})
	message := &mailer.Message{
		From: mail.Address{
			Name:    senderName,
			Address: senderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Opt-In",
		Text:    text,
		HTML:    html,
	}
	return mailClient.Send(message)
}

func addEmailToMailingList(db *daos.Dao, email string) error {
	mailinglist, err := db.FindCollectionByNameOrId("mailinglist")
	if err != nil {
		return err
	}
	existing, err := db.FindFirstRecordByData("mailinglist", "email", email)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}
	record := models.NewRecord(mailinglist)
	record.RefreshId()
	record.Set("email", email)
	return db.SaveRecord(record)
}

func removeEmailToMailingList(db *daos.Dao, email string) error {
	record, err := db.FindFirstRecordByData("mailinglist", "email", email)
	if err != nil {
		return err
	}
	return db.DeleteRecord(record)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && !strings.Contains(email, " ")
}

type templateData struct {
	Name    string
	Email   string
	Message string
	Link    string
}

func getTemplates(template string, data templateData) (string, string, error) {
	tt, err := ttemplate.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.txt.gotpl", template)))
	if err != nil {
		return "", "", err
	}

	ht, err := htemplate.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.html.gotpl", template)))
	if err != nil {
		return "", "", err
	}

	twr := &bytes.Buffer{}
	if err = tt.Execute(twr, data); err != nil {
		return "", "", err
	}

	hwr := &bytes.Buffer{}
	if err = ht.Execute(hwr, data); err != nil {
		return "", "", err
	}
	return twr.String(), hwr.String(), nil
}
