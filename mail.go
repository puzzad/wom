package wom

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	htemplate "html/template"
	"net/mail"
	"strings"
)

//go:embed templates
var templates embed.FS
var templateFS = echo.MustSubFS(templates, "templates")

func SendContactFormMail(mailClient mailer.Mailer, contactEmail, senderName, senderAddress, email string, name string, content string) error {
	html, err := GetTemplates("contact", TemplateData{
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
		HTML:    html,
	}
	return mailClient.Send(message)
}

func sendSubscriptionConfirmedMail(mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createUnsubscribeJwt(mailingListSecret, email)
	if err != nil {
		return err
	}
	html, err := GetTemplates("subscribed", TemplateData{
		SiteURL: siteURL,
		Token:   token,
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
		HTML:    html,
	}
	return mailClient.Send(message)
}

func sendSubscriptionUnsubscribedMail(mailClient mailer.Mailer, siteURL, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createSubscriptionJwt(mailingListSecret, email)
	if err != nil {
		fmt.Printf("Error creating subscription JWT: %s", err)
		return err
	}
	html, err := GetTemplates("unsubscribed", TemplateData{
		SiteURL: siteURL,
		Token:   token,
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
	html, err := GetTemplates("optin", TemplateData{
		Token:   token,
		SiteURL: siteURL,
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
		HTML:    html,
	}
	return mailClient.Send(message)
}

func addEmailToMailingList(db *daos.Dao, email string) error {
	mailinglist, err := db.FindCollectionByNameOrId("mailinglist")
	if err != nil {
		return err
	}
	existing, _ := db.FindFirstRecordByData("mailinglist", "email", email)
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

type TemplateData struct {
	Name    string
	Email   string
	Message string
	Token   string
	SiteURL string
}

func GetTemplates(template string, data TemplateData) (string, error) {
	templateName := fmt.Sprintf("template-%s.gotpl", template)
	parsedTemplates, err := htemplate.ParseFS(templateFS, "common-*.gotpl", templateName)
	if err != nil {
		return "", err
	}
	hwr := &bytes.Buffer{}
	err = parsedTemplates.ExecuteTemplate(hwr, templateName, data)
	if err != nil {
		return "", err
	}
	return hwr.String(), nil
}
