package wom

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewImportCmd(app *pocketbase.PocketBase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports adventures",
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("filename")
			apiurl, _ := cmd.Flags().GetString("apiurl")
			email, _ := cmd.Flags().GetString("email")
			password, _ := cmd.Flags().GetString("password")
			cobra.CheckErr(uploadAdventures(filename, apiurl, email, password))
		},
	}
	cmd.Flags().String("filename", "adventures.zip", "path to an adventures zip file")
	cmd.Flags().Bool("production", false, "Whether this is an upload to production or dev")
	cmd.Flags().String("apiurl", "http://localhost:8090/", "URL of the WOM instance")
	cmd.Flags().String("email", "", "Admin email address")
	cmd.Flags().String("password", "", "Admin password")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func updateAdventures(app *pocketbase.PocketBase, adventures []*adventure) {
	puzzlesfs, err := app.NewFilesystem()
	defer func() {
		_ = puzzlesfs.Close()
	}()
	if err != nil {
		log.Fatalf("Unable to get filesystem: %s", err)
	}
	for i := range adventures {
		adventureRecord, err := app.Dao().FindFirstRecordByData("adventures", "name", adventures[i].name)
		if err != nil {
			collection, err := app.Dao().FindCollectionByNameOrId("adventures")
			if err != nil {
				log.Fatalf("Unable to find adventures collection: %s\n", err)
			}
			adventureRecord = models.NewRecord(collection)
			adventureRecord.RefreshId()
		}
		adventureRecord.Set("name", adventures[i].name)
		adventureRecord.Set("description", adventures[i].description)
		adventureRecord.Set("price", adventures[i].price)
		adventureRecord.Set("public", !adventures[i].private)
		adventureRecord.Set("features", adventures[i].features)
		if adventures[i].background != nil {
			adventureRecord.Set("background", "background.jpg")
			err = puzzlesfs.Upload(adventures[i].background, adventureRecord.BaseFilesPath()+"/background.jpg")
			if err != nil {
				log.Fatalf("Unable to upload background: %s", err)
			}
		}
		if adventures[i].logo != nil {
			adventureRecord.Set("logo", "logo.png")
			err = puzzlesfs.Upload(adventures[i].logo, adventureRecord.BaseFilesPath()+"/logo.png")
			if err != nil {
				log.Fatalf("Unable to upload logo: %s", err)
			}
		}
		err = app.Dao().SaveRecord(adventureRecord)
		if err != nil {
			log.Fatalf("Unable to save adventure: %s\n", err)
		}
		for j := range adventures[i].puzzles {
			puzzleRecords, err := app.Dao().FindRecordsByExpr("puzzles", dbx.HashExp{"title": adventures[i].puzzles[j].name, "adventure": adventureRecord.Id})
			if err != nil {
				log.Fatalf("Unable to find puzzle: %s\n", err)
			}
			var puzzleRecord *models.Record
			if len(puzzleRecords) == 0 {
				collection, err := app.Dao().FindCollectionByNameOrId("puzzles")
				if err != nil {
					log.Fatalf("Unable to find puzzles collection: %s\n", err)
				}
				puzzleRecord = models.NewRecord(collection)
				puzzleRecord.RefreshId()
			} else {
				puzzleRecord = puzzleRecords[0]
			}
			puzzleRecord.Set("title", adventures[i].puzzles[j].name)
			puzzleRecord.Set("adventure", adventureRecord.Id)
			puzzleRecord.Set("information", adventures[i].puzzles[j].info)
			puzzleRecord.Set("story", adventures[i].puzzles[j].story)
			puzzleRecord.Set("puzzle", adventures[i].puzzles[j].text)
			err = app.Dao().SaveRecord(puzzleRecord)
			if err != nil {
				log.Fatalf("Unable to save puzzle: %s\n", err)
			}
			for k := range adventures[i].puzzles[j].answers {
				answerRecords, err := app.Dao().FindRecordsByExpr("answers", dbx.HashExp{"puzzle": puzzleRecord.Id, "content": adventures[i].puzzles[j].answers[k]})
				if err != nil {
					log.Fatalf("Unable to find answer: %s\n", err)
				}
				var answerRecord *models.Record
				if len(answerRecords) == 0 {
					collection, err := app.Dao().FindCollectionByNameOrId("answers")
					if err != nil {
						log.Fatalf("Unable to find answers collection: %s\n", err)
					}
					answerRecord = models.NewRecord(collection)
					answerRecord.RefreshId()
				} else {
					answerRecord = answerRecords[0]
				}
				answerRecord.Set("puzzle", puzzleRecord.Id)
				answerRecord.Set("content", adventures[i].puzzles[j].answers[k])
				err = app.Dao().SaveRecord(answerRecord)
				if err != nil {
					log.Fatalf("Unable to save answer: %s\n", err)
				}
			}
			for k := range adventures[i].puzzles[j].hints {
				hintRecords, err := app.Dao().FindRecordsByExpr("hints", dbx.HashExp{
					"puzzle":  puzzleRecord.Id,
					"title":   adventures[i].puzzles[j].hints[k][0],
					"message": adventures[i].puzzles[j].hints[k][1],
				})
				if err != nil {
					log.Fatalf("Unable to find hint: %s\n", err)
				}
				var hintRecord *models.Record
				if len(hintRecords) == 0 {
					collection, err := app.Dao().FindCollectionByNameOrId("hints")
					if err != nil {
						log.Fatalf("Unable to find hints collection: %s\n", err)
					}
					hintRecord = models.NewRecord(collection)
					hintRecord.RefreshId()
				} else {
					hintRecord = hintRecords[0]
				}
				hintRecord.Set("puzzle", puzzleRecord.Id)
				hintRecord.Set("title", adventures[i].puzzles[j].hints[k][0])
				hintRecord.Set("message", adventures[i].puzzles[j].hints[k][1])
				hintRecord.Set("order", k)
				err = app.Dao().SaveRecord(hintRecord)
				if err != nil {
					log.Fatalf("Unable to save hints: %s\n", err)
				}
			}
			if j == 0 {
				adventureRecord.Set("firstpuzzle", puzzleRecord.Id)
				err = app.Dao().SaveRecord(adventureRecord)
				if err != nil {
					log.Fatalf("Unable to update firstPuzzle: %s\n", err)
				}
			}
		}
		allPuzzlesSlice, err := app.Dao().FindRecordsByExpr("puzzles", dbx.HashExp{"adventure": adventureRecord.Id})
		if err != nil {
			log.Fatalf("Unable to update puzzle order: %s", err)
		}
		idMap := make(map[string]string)
		nextMap := make(map[string]string)
		for i := range allPuzzlesSlice {
			idMap[allPuzzlesSlice[i].Get("title").(string)] = allPuzzlesSlice[i].Id
		}
		for j := range adventures[i].puzzles {
			if j > 0 {
				nextMap[adventures[i].puzzles[j-1].name] = idMap[adventures[i].puzzles[j].name]
			}
		}
		for j := range allPuzzlesSlice {
			allPuzzlesSlice[j].Set("next", nextMap[allPuzzlesSlice[j].Get("title").(string)])
			err = app.Dao().SaveRecord(allPuzzlesSlice[j])
			if err != nil {
				log.Fatalf("Unable to update puzzle order: %s", err)
			}
		}
	}
}

func uploadAdventures(fileName, apiurl string, email string, password string) error {
	data := struct {
		Token string
	}{}
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open zip: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()
	//TODO: Should remove resty really, its overkill now I'm just uploading a file
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(3*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetBaseURL(apiurl).
		SetHeader("Content-Type", "application/json")
	resp, err := client.R().
		SetBody(`{"identity":"` + email + `", "password":"` + password + `"}`).
		SetResult(&data).
		Post("api/admins/auth-with-password")
	if err != nil {
		return fmt.Errorf("unable to login")
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to login: Status %d", resp.StatusCode())
	}
	client.SetAuthToken(data.Token)
	resp, err = client.R().
		SetContentLength(true).
		SetFileReader("adventures.zip", "adventures.zip", file).
		Post("/import/zip")
	if err != nil {
		return fmt.Errorf("unable to upload zip: %s", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unable to upload zip: Status %d", resp.StatusCode())
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

	if len(puzzles) == 0 {
		return nil
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
	content := replaceVariables(zfs, readTextFile(zfs, filepath.Join(folder, "puzzle.html")), folder)
	story := replaceVariables(zfs, readTextFile(zfs, filepath.Join(folder, "story.html")), folder)
	information := replaceVariables(zfs, readTextFile(zfs, filepath.Join(folder, "information.html")), folder)

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
	} else {
		output = content
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
