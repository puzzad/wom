package wom

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func NewImportCmd(app *pocketbase.PocketBase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports adventures",
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("filename")
			_ = zippy(filename)
		},
	}
	cmd.Flags().String("filename", "adventures.zip", "path to an adventures zip file")
	cmd.Flags().Bool("production", false, "Whether this is an upload to production or dev")
	return cmd
}

func zippy(fileName string) error {
	archive, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	adventures := getAdventures(archive, false)
	uploadAdventures(archive, adventures)
	return nil
}

func uploadAdventures(zfs fs.FS, adventures []*adventure) {
	for i := range adventures {
		updatePuzzleOrder(adventures[i].puzzles)
		fmt.Printf("%s\n", adventures[i].name)
		for j := range adventures[i].puzzles {
			fmt.Printf("\t%s\n", adventures[i].puzzles[j].name)
		}
	}
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
