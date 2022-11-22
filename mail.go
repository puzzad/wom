package wom

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	htemplate "html/template"
	"path/filepath"
	ttemplate "text/template"

	"github.com/mailgun/mailgun-go/v4"
)

var (
	mailDomain = flag.String("mailgun-domain", "", "Domain to use when sending mail from mailgun")
	mailSender = flag.String("mailgun-sender", "", "From e-mail address to use for e-mails")
	apiKey     = flag.String("mailgun-api-key", "", "API key to use for mailgun")
	apiBase    = flag.String("mailgun-api-base", "https://api.eu.mailgun.net/v3", "Base URL for the mailgun API")

	subscriptionConfirmLink     = flag.String("subscription-confirm-link", "", "Link to send users to confirm their subscription to the mailing list")
	subscriptionUnsubscribeLink = flag.String("subscription-unsubscribe-link", "", "Link to send users to unsubscribe from the mailing list")
)

func sendSubscriptionOptInMail(ctx context.Context, email string) error {
	token, err := createSubscriptionJwt(email)
	if err != nil {
		return err
	}

	return sendMail(ctx, email, "Puzzad: mailing list opt-in", "optin", map[string]string{
		"Link": fmt.Sprintf(*subscriptionConfirmLink, token),
	})
}

func sendSubscriptionConfirmedMail(ctx context.Context, email string) error {
	token, err := createUnsubscribeJwt(email)
	if err != nil {
		return err
	}

	return sendMail(ctx, email, "Puzzad: mailing list confirmation", "subscribed", map[string]string{
		"Unsubscribe": fmt.Sprintf(*subscriptionUnsubscribeLink, token),
	})
}

func sendSubscriptionEndedMail(ctx context.Context, email string) error {
	token, err := createSubscriptionJwt(email)
	if err != nil {
		return err
	}

	return sendMail(ctx, email, "Puzzad: unsubscribed", "unsubscribed", map[string]string{
		"Link": fmt.Sprintf(*subscriptionConfirmLink, token),
	})
}

func sendMail(ctx context.Context, address, subject, template string, data any) error {
	mg := mailgun.NewMailgun(*mailDomain, *apiKey)
	mg.SetAPIBase(*apiBase)

	tt, err := ttemplate.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.txt.gotpl", template)))
	if err != nil {
		return err
	}

	ht, err := htemplate.ParseFiles(filepath.Join("templates", fmt.Sprintf("%s.html.gotpl", template)))
	if err != nil {
		return err
	}

	twr := &bytes.Buffer{}
	if err = tt.Execute(twr, data); err != nil {
		return err
	}

	hwr := &bytes.Buffer{}
	if err = ht.Execute(hwr, data); err != nil {
		return err
	}

	message := mg.NewMessage(*mailSender, subject, twr.String(), address)
	message.SetHtml(hwr.String())
	_, _, err = mg.Send(ctx, message)
	return err
}
