package wom

import (
	"fmt"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"net/mail"
	"strings"
)

func SendContactFormMail(mailClient mailer.Mailer, contactEmail, senderName, senderAddress, email string, name string, content string) error {
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
		Text:    fmt.Sprintf("Name: %s\nMessage\n%s", name, content),
	}
	return mailClient.Send(message)
}

func sendSubscriptionConfirmedMail(mailClient mailer.Mailer, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createUnsubscribeJwt(mailingListSecret, email)
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
		Text:    fmt.Sprintf("/mail/unsubscribe/%s", token),
	}
	return mailClient.Send(message)
}

func sendSubscriptionUnsubscribedMail(mailClient mailer.Mailer, senderName, senderAddress, email string) error {
	message := &mailer.Message{
		From: mail.Address{
			Name:    senderName,
			Address: senderAddress,
		},
		To: []mail.Address{{
			Address: email,
		}},
		Subject: "Mailinglist Unsubscribed",
		Text:    fmt.Sprintf("Sorry to see you go"),
	}
	return mailClient.Send(message)
}

func sendSubscriptionOptInMail(mailClient mailer.Mailer, senderName, senderAddress, mailingListSecret, email string) error {
	token, err := createSubscriptionJwt(mailingListSecret, email)
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
		Subject: "Mailinglist Opt-In",
		Text:    fmt.Sprintf("/mail/confirm/%s", token),
	}
	return mailClient.Send(message)
}

func addEmailToMailingList(db *daos.Dao, email string) error {
	mailinglist, err := db.FindCollectionByNameOrId("mailinglist")
	if err != nil {
		return err
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
