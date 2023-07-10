package wom

import (
	"archive/zip"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"io"
	"net/http"
)

func createWomRoutesHook(app *pocketbase.PocketBase) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		_, err := e.Router.AddRoute(echo.Route{
			Name:    "start adventure",
			Path:    "/adventure/:id/start",
			Method:  http.MethodPost,
			Handler: startAdventure(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Name:    "start game",
			Path:    "/games/:code/start",
			Method:  http.MethodPost,
			Handler: startGame(app),
		})
		if err != nil {
			return err
		}
		_, err = e.Router.AddRoute(echo.Route{
			Name:    "send contact form",
			Path:    "/mail/contact",
			Method:  http.MethodPost,
			Handler: sendContactForm(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:      http.MethodPost,
			Name:        "import adventure zip",
			Path:        "/import/zip",
			Middlewares: []echo.MiddlewareFunc{apis.RequireAdminAuth()},
			Handler:     importAdventures(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Name:    "request hint",
			Path:    "/hints/request",
			Handler: requestHint(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Name:    "request hint",
			Path:    "/mail/subscribe",
			Handler: handleSubscribe(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Name:    "request hint",
			Path:    "/mail/confirm",
			Handler: handleConfirm(app),
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Name:    "request hint",
			Path:    "/mail/unsubscribe",
			Handler: handleUnsubscribe(app),
		})
		return nil
	}
}

func handleUnsubscribe(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		secretKey, _ := app.RootCmd.Flags().GetString("mailinglistSecretKey")
		email, err := validateSubscriptionJWT(secretKey, "unsubscribe", data.Token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err := removeEmailToMailingList(app, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionUnsubscribedMail(app, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleConfirm(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		type req struct {
			Token string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		secretKey, _ := app.RootCmd.Flags().GetString("mailinglistSecretKey")
		email, err := validateSubscriptionJWT(secretKey, "unsubscribe", data.Token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		if err := removeEmailToMailingList(app, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionConfirmedMail(app, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func handleSubscribe(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		authInfo, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		type req struct {
			Email   string
			Captcha string
		}
		var data = req{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}
		email := ""
		if authInfo != nil && authInfo.Verified() {
			email = authInfo.Email()
		}
		if data.Email != email {
			secretKey, _ := app.RootCmd.Flags().GetString("hcaptchaSecretKey")
			siteKey, _ := app.RootCmd.Flags().GetString("hcaptchaSiteKey")
			if err := checkCaptcha(siteKey, secretKey, data.Captcha); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			}
			if err := sendSubscriptionOptInMail(app, data.Email); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
			}
			return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": true})
		}
		if err := addEmailToMailingList(app, data.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		if err := sendSubscriptionConfirmedMail(app, data.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"NeedConfirm": false})
	}
}

func requestHint(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		game, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if game == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Game not found"})
		}
		usedhints, err := app.Dao().FindCollectionByNameOrId("usedhints")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "usedhints not found"})
		}
		guesses, err := app.Dao().FindCollectionByNameOrId("guesses")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "guesses not found"})
		}
		record := models.NewRecord(usedhints)
		record.RefreshId()
		hintID, err := io.ReadAll(c.Request().Body)
		defer func() {
			_ = c.Request().Body.Close()
		}()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "hintID not found"})
		}
		record.Set("hint", hintID)
		record.Set("game", game.Id)
		err = app.Dao().SaveRecord(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to request hint"})
		}
		record = models.NewRecord(guesses)
		record.RefreshId()
		record.Set("game", game.Id)
		record.Set("puzzle", game.Get("puzzle"))
		record.Set("content", "*hint")
		err = app.Dao().SaveRecord(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to request hint"})
		}
		return c.JSON(http.StatusOK, "")
	}
}

func importAdventures(app *pocketbase.PocketBase) func(c echo.Context) error {
	return func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		if len(form.File["adventures.zip"]) != 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Must be one file called adventures.zip"})
		}
		file := form.File["adventures.zip"][0]
		fileReader, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to open file"})
		}
		zipReader, err := zip.NewReader(fileReader, file.Size)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to open zip"})
		}
		adventures := getAdventures(zipReader, false)
		updateAdventures(app, adventures)
		return c.JSON(http.StatusOK, nil)
	}
}
