package wom

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/render"
)

func SubscribeToMailingList(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email   string
		Captcha string
	}

	type res struct {
		NeedConfirm bool
	}

	var data = req{}
	if err := render.DecodeJSON(r.Body, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// No point bothering if the e-mail address isn't right
	if !validEmail(data.Email) {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid e-mail address")))
		return
	}

	// If they're not using a verified e-mail address, we should expect a captcha and then send an opt-in
	if data.Email != getEmailFromJwt(r) {
		if err := checkCaptcha(data.Captcha); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if err := sendSubscriptionOptInMail(r.Context(), data.Email); err != nil {
			log.Printf("Unable to send email: %v", err)
			render.Render(w, r, ErrInternalError(err))
			return
		}

		render.JSON(w, r, res{NeedConfirm: true})
		return
	}

	// They're using an e-mail address they've previously validated, just sign them up
	if err := addEmailToMailingList(r.Context(), data.Email); err != nil {
		log.Printf("Unable to add user to mailing list: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	if err := sendSubscriptionConfirmedMail(r.Context(), data.Email); err != nil {
		log.Printf("Unable to send email: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	render.JSON(w, r, res{NeedConfirm: false})
}

func ConfirmMailingListSubscription(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Token string
	}

	var data = req{}
	if err := render.DecodeJSON(r.Body, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	email, err := validateSubscriptionJwt("subscribe", data.Token)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := addEmailToMailingList(r.Context(), email); err != nil {
		log.Printf("Unable to add user to mailing list: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	if err := sendSubscriptionConfirmedMail(r.Context(), email); err != nil {
		log.Printf("Unable to send email: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UnsubscribeFromMailingList(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Token string
	}

	var data = req{}
	if err := render.DecodeJSON(r.Body, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	email, err := validateSubscriptionJwt("unsubscribe", data.Token)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := removeEmailFromMailingList(r.Context(), email); err != nil {
		log.Printf("Unable to remove user from mailing list: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	if err := sendSubscriptionEndedMail(r.Context(), email); err != nil {
		log.Printf("Unable to send email: %v", err)
		render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func SendContactForm(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Token   string
		Name    string
		Email   string
		Message string
	}
	var data = req{}
	if err := render.DecodeJSON(r.Body, &data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if strings.TrimSpace(data.Name) == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("name is required")))
		return
	}
	if strings.TrimSpace(data.Message) == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("message is required")))
		return
	}

	if data.Email != getEmailFromJwt(r) {
		if err := checkCaptcha(data.Token); err != nil {
			log.Printf("%s", err.Error())
			_ = render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		if !validEmail(data.Email) {
			_ = render.Render(w, r, ErrInvalidRequest(errors.New("invalid email")))
			return
		}
	}

	if err := SendContactFormMail(r.Context(), data.Email, data.Name, data.Message); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && !strings.Contains(email, " ")
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func ErrInternalError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal error",
		ErrorText:      err.Error(),
	}
}
