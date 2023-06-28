package main

import (
	"flag"
	"fmt"
	"github.com/csmith/aca"
	"github.com/csmith/envflag"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/puzzad/wom"
	_ "github.com/puzzad/wom/migrations"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var port = flag.Int("port", 3000, "Port to listen for HTTP requests")
var initialAdminEmail = flag.String("au", "pocketbase@greg.holmes.name", "")
var initialAdminPassword = flag.String("ap", "password", "")
var dataDir = flag.String("dir", "./data", "")

func main() {
	envflag.Parse()
	app := pocketbase.NewWithConfig(&pocketbase.Config{
		DefaultDataDir: *dataDir,
	})
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		if len(*initialAdminEmail) == 0 ||
			is.EmailFormat.Validate(*initialAdminEmail) != nil ||
			len(*initialAdminPassword) == 0 {
			return fmt.Errorf("invalid admin credentials")
		}
		admin, err := e.App.Dao().FindAdminByEmail(*initialAdminEmail)
		if err != nil {
			admin = &models.Admin{
				Email: *initialAdminEmail,
			}
		}
		err = admin.SetPassword(*initialAdminPassword)
		if err != nil {
			return err
		}
		return e.App.Dao().SaveAdmin(admin)
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		_, err := e.Router.AddRoute(echo.Route{
			Name:    "start adventure",
			Path:    "/adventure/:id/start",
			Method:  http.MethodPost,
			Handler: startAdventure(app),
		})
		if err != nil {
			return err
		}
		_, err = e.Router.AddRoute(echo.Route{
			Name:    "send contact form",
			Path:    "/mail/contact",
			Method:  http.MethodPost,
			Handler: wom.SendContactForm,
		})
		//e.Router.Add(http.MethodGet, "/mail/subscribe", wom.SubscribeToMailingList)
		//e.Router.Add(http.MethodGet, "/mail/confirm", wom.ConfirmMailingListSubscription)
		//e.Router.Add(http.MethodGet, "/mail/unsubscribe", wom.UnsubscribeFromMailingList)
		//e.Router.Add(http.MethodGet, "/mail/contact", wom.SendContactForm)
		return nil
	})
	app.OnRecordBeforeUpdateRequest("adventures").Add(func(e *core.RecordUpdateEvent) error {
		return preserveOriginalFilenames(e.UploadedFiles, e.Record)
	})
	app.OnRecordBeforeCreateRequest("adventures").Add(func(e *core.RecordCreateEvent) error {
		return preserveOriginalFilenames(e.UploadedFiles, nil)
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func startAdventure(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")
		if len(id) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Adventure"})
		}
		adventure, err := app.Dao().FindRecordById("adventures", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Adventure not found"})
		}
		collection, err := app.Dao().FindCollectionByNameOrId("games")
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Adventure not found"})
		}
		user, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if user == nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
		}
		if !user.Verified() {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email must be verified"})
		}
		acaGen, err := aca.NewGenerator("-", rand.NewSource(time.Now().UnixNano()))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to generate ACA"})
		}
		code := acaGen.Generate()
		record := models.NewRecord(collection)
		form := forms.NewRecordUpsert(app, record)
		err = form.LoadData(map[string]any{
			"status":    "PAID",
			"user":      user.Id,
			"adventure": adventure.Id,
			"code":      code,
		})
		if err = form.Submit(); err != nil {
			fmt.Printf("%v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		if err = form.Submit(); err != nil {
			fmt.Printf("%v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to add adventure"})
		}
		return c.JSON(http.StatusOK, map[string]string{"code": code})
	}
}

func preserveOriginalFilenames(uploadedFiles map[string][]*filesystem.File, record *models.Record) error {
	oldNames := map[string]map[string]string{}
	for field, files := range uploadedFiles {
		if len(files) == 0 {
			continue
		}
		oldNames[field] = make(map[string]string, len(files))
		for _, f := range files {
			oldNames[field][f.Name] = f.OriginalName
			f.Name = f.OriginalName
		}
	}
	for field, filenames := range oldNames {
		files := record.GetStringSlice(field)

		for i, old := range files {
			if newName, ok := filenames[old]; ok {
				files[i] = newName
			}
		}
		record.Set(field, files)
	}
	return nil
}
