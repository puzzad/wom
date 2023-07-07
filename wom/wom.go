package wom

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/csmith/aca"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConfigurePocketBase(app *pocketbase.PocketBase) {
	serveCmd := cmd.NewServeCommand(app, true)
	serveCmd.PersistentFlags().StringP("email", "e", "", "Sets the initial admin email")
	serveCmd.PersistentFlags().StringP("password", "p", "", "Sets the initial admin password")
	serveCmd.PersistentFlags().StringP("webhook-url", "w", "", "Webhook to send events to {'content': 'message'}")
	app.RootCmd.AddCommand(serveCmd)
	app.RootCmd.AddCommand(NewImportCmd(app))
	app.OnBeforeServe().Add(createAdminAccountHook(serveCmd))
	app.OnBeforeServe().Add(createWomRoutesHook(app))
	app.OnRecordBeforeUpdateRequest("adventures").Add(createPreserveFilenameUpdateHook)
	app.OnRecordBeforeCreateRequest("adventures").Add(createPreserveFilenameCreateHook)
	app.OnRecordBeforeCreateRequest("guesses").Add(createBeforeGuessCreatedHook(app))
	app.OnRecordAfterCreateRequest("guesses").Add(createGuessCreatedHook(app))
}

func createBeforeGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		e.Record.Set("correct", checkGuess(app, e.Record))
		return nil
	}
}

func createGuessCreatedHook(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		var code, title string
		err := app.Dao().DB().Select("games.username as username", "puzzles.title as title").
			From("guesses").
			InnerJoin("games", dbx.NewExp("games.id=guesses.game")).
			InnerJoin("puzzles", dbx.NewExp("puzzles.id=guesses.puzzle")).
			Where(dbx.HashExp{"guesses.id": e.Record.Id}).
			Row(&code, &title)
		if err == nil {
			webhookURL, _ := app.RootCmd.Flags().GetString("webhook-url")
			if e.Record.Get("correct").(bool) {
				sendWebhook(webhookURL, fmt.Sprintf(":tada: %s/%s: %s", code, title, e.Record.Get("content")))
			} else {
				sendWebhook(webhookURL, fmt.Sprintf(":x: %s/%s: %s", code, title, e.Record.Get("content")))
			}
		}
		return nil
	}
}

func createPreserveFilenameCreateHook(e *core.RecordCreateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
}

func createPreserveFilenameUpdateHook(e *core.RecordUpdateEvent) error {
	return preserveOriginalFilenames(e.UploadedFiles, e.Record)
}

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
			Name:   "send contact form",
			Path:   "/mail/contact",
			Method: http.MethodPost,
			//Handler: wom.SendContactForm,
		})
		_, err = e.Router.AddRoute(echo.Route{
			Method:      http.MethodPost,
			Name:        "import adventure zip",
			Path:        "/import/zip",
			Middlewares: []echo.MiddlewareFunc{apis.RequireAdminAuth()},
			Handler:     importAdventures,
		})
		//e.Router.Add(http.MethodGet, "/mail/subscribe", wom.SubscribeToMailingList)
		//e.Router.Add(http.MethodGet, "/mail/confirm", wom.ConfirmMailingListSubscription)
		//e.Router.Add(http.MethodGet, "/mail/unsubscribe", wom.UnsubscribeFromMailingList)
		//e.Router.Add(http.MethodGet, "/mail/contact", wom.SendContactForm)
		return nil
	}
}

func importAdventures(c echo.Context) error {
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
	adventureMap := make(map[string][]string)
	for i := range adventures {
		puzzleList := make([]string, 0)
		for j := range adventures[i].puzzles {
			puzzleList = append(puzzleList, adventures[i].puzzles[j].name)
		}
		adventureMap[adventures[i].name] = puzzleList
	}
	fmt.Printf("Uploaded: %+v", adventureMap)
	return c.JSON(http.StatusOK, adventureMap)
}

func createAdminAccountHook(serveCmd *cobra.Command) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		email, _ := serveCmd.Flags().GetString("email")
		password, _ := serveCmd.Flags().GetString("password")
		if len(email) == 0 && len(password) == 0 {
			return nil
		}
		if is.EmailFormat.Validate(email) != nil || len(password) <= 5 {
			return fmt.Errorf("invalid admin credentials\n")
		}
		admin, err := e.App.Dao().FindAdminByEmail(email)
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
		err = e.App.Dao().SaveAdmin(admin)
		if err != nil {
			fmt.Printf("Error saving admin: %v\n", err)
		}
		return err
	}
}

func sendWebhook(webHookURL, message string) {
	type webhook struct {
		Content string `json:"content"`
	}
	if len(webHookURL) == 0 {
		return
	}
	data := &webhook{
		Content: message,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, err = http.Post(webHookURL, "application/json", bytes.NewReader(dataBytes))
	if err != nil {
		return
	}
}

func checkGuess(app *pocketbase.PocketBase, r *models.Record) bool {
	content := r.Get("content")
	puzzle := r.Get("puzzle")
	game := r.Get("game")
	var count string
	err := app.Dao().DB().Select("count(*)").From("answers").
		AndWhere(dbx.HashExp{"puzzle": puzzle}).
		AndWhere(dbx.HashExp{"game": game}).
		AndWhere(dbx.HashExp{"answer": content}).
		Row(&count)
	if err != nil {
		return false
	}
	return count == "1"
}

func startGame(app *pocketbase.PocketBase) func(echo.Context) error {
	return func(c echo.Context) error {
		code := c.PathParam("code")
		if len(code) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Game"})
		}
		q := app.Dao().DB().NewQuery("UPDATE games SET status = 'ACTIVE', puzzle = (SELECT adventures.firstpuzzle FROM adventures WHERE adventures.id = games.adventure), start = datetime('now') WHERE code = {:code} AND status = 'PAID' AND (puzzle='' OR puzzle IS NULL);")
		q = q.Bind(dbx.Params{"code": code})
		result, err := q.Execute()
		if err != nil {
			fmt.Printf("%v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 1"})
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 2"})
		}
		if rows != 1 {
			fmt.Printf("Rows: %d\n", rows)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to start game 3"})
		}
		return nil
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

func getAdventures(zfs fs.FS, prod bool) []*adventure {
	var adventures []*adventure
	for _, f := range readDir(zfs, ".") {
		if f.IsDir() {
			checkedAdventure := checkAdventure(zfs, f.Name(), prod)
			if checkedAdventure != nil {
				adventures = append(adventures, checkedAdventure)
			}
		}
	}
	return adventures
}

func checkAdventure(zfs fs.FS, name string, prod bool) *adventure {
	if !exists(zfs, filepath.Join(name, "PRICE")) {
		return nil
	}

	price := readTextFile(zfs, filepath.Join(name, "PRICE"))
	description := readTextFile(zfs, filepath.Join(name, "description.html"))
	private := exists(zfs, filepath.Join(name, "PRIVATE"))
	devOnly := exists(zfs, filepath.Join(name, "DEVONLY"))
	features := parseJsonFile(zfs, filepath.Join(name, "features.json"), &adventureFeatures{})

	if devOnly && prod {
		return nil
	}

	var puzzles []puzzle
	for _, f := range readDir(zfs, name) {
		if f.IsDir() {
			puzzles = append(puzzles, checkPuzzle(zfs, filepath.Join(name, f.Name())))
		}
	}

	return &adventure{
		name:        name,
		description: description,
		price:       price,
		private:     private,
		dev:         devOnly,
		puzzles:     puzzles,
		background:  readBinaryFile(zfs, filepath.Join(name, "background.jpg")),
		logo:        readBinaryFile(zfs, filepath.Join(name, "logo.png")),
		preview:     readBinaryFile(zfs, filepath.Join(name, "preview.png")),
		features:    features,
	}
}

func checkPuzzle(zfs fs.FS, folder string) puzzle {
	parts := strings.SplitN(filepath.Base(folder), " - ", 2)
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Puzzle folder has invalid number: %s (%v)", folder, err)
	}
	name := parts[1]
	answers := readTextFileLines(zfs, filepath.Join(folder, "answers.txt"))
	var hints [][2]string
	for _, h := range readTextFileLines(zfs, filepath.Join(folder, "hints.txt")) {
		hintParts := strings.SplitN(h, ": ", 2)
		hints = append(hints, [2]string{hintParts[0], hintParts[1]})
	}
	content := replaceVariables(zfs, folder, readTextFile(zfs, filepath.Join(folder, "puzzle.html")))
	story := replaceVariables(zfs, folder, readTextFile(zfs, filepath.Join(folder, "information.html")))
	information := replaceVariables(zfs, folder, readTextFile(zfs, filepath.Join(folder, "story.html")))

	return puzzle{
		order:   number,
		name:    name,
		slug:    fmt.Sprintf("%s%0.2d", strings.ToLower(filepath.Dir(folder)), number),
		text:    content,
		story:   story,
		info:    information,
		answers: answers,
		hints:   hints,
		files:   nil,
	}
}

func updatePuzzleOrder(puzzles []puzzle) {
	for i := range puzzles {
		if i > 0 {
			puzzles[i-1].next = puzzles[i].slug
		}
	}
}

func parseJsonFile[T interface{}](zfs fs.FS, name string, parsed *T) *T {
	dataString := readTextFile(zfs, name)
	err := json.Unmarshal([]byte(dataString), parsed)
	if err != nil {
		log.Fatalf("Unable to parse features: %s", err)
	}
	return parsed
}

func replaceVariables(zfs fs.FS, content, path string) (output string) {
	matcher := regexp.MustCompile(`\$(.*?)\$`)
	matches := matcher.FindAllString(content, -1)
	if matches != nil {
		for _, match := range matches {
			filename := match[1 : len(match)-1]
			if len(filename) == 0 {
				log.Fatalf("Invalid variable in file: %s: %s", match, path)
			}
			fileBytes := readBinaryFile(zfs, filepath.Join(path, filename))
			contentType := http.DetectContentType(fileBytes)
			b64 := base64.StdEncoding.EncodeToString(fileBytes)
			datauri := fmt.Sprintf("data:%s;base64,%s", contentType, b64)
			output = strings.ReplaceAll(content, match, datauri)
		}
	}
	return output
}

func readBinaryFile(zfs fs.FS, name string) []byte {
	b, err := fs.ReadFile(zfs, name)
	if err != nil {
		log.Fatalf("Unable to read file %s: %v", name, err)
	}
	return b
}

func readTextFile(zfs fs.FS, name string) string {
	return string(readBinaryFile(zfs, name))
}

func readTextFileLines(zfs fs.FS, name string) []string {
	return strings.Split(strings.TrimSpace(readTextFile(zfs, name)), "\n")
}

func exists(zfs fs.FS, name string) bool {
	_, err := fs.Stat(zfs, name)
	return err == nil
}

func readDir(zfs fs.FS, name string) []os.DirEntry {
	files, err := fs.ReadDir(zfs, name)
	if err != nil {
		panic(err)
	}
	return files
}

type adventure struct {
	name        string
	description string
	price       string
	private     bool
	dev         bool
	background  []byte
	logo        []byte
	preview     []byte
	features    *adventureFeatures
	puzzles     []puzzle
}

type adventureFeatures struct {
	Difficulty    string                 `json:"difficulty"`
	Players       string                 `json:"players"`
	Equipment     string                 `json:"equipment"`
	Puzzles       string                 `json:"puzzles"`
	Accessibility adventureAccessibility `json:"accessibility"`
}

type adventureAccessibility struct {
	Hearing string `json:"hearing"`
	Vision  string `json:"vision"`
	Colours string `json:"colours"`
	Motion  string `json:"motion"`
}

type puzzle struct {
	order   int
	name    string
	slug    string
	text    string
	info    string
	story   string
	answers []string
	hints   [][2]string
	files   []string
	next    string
}
