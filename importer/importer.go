package importer

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func uploadAdventures(adventures []*adventure) {
	for i := range adventures {
		fmt.Printf("Adventure: %s\n", adventures[i].name)
	}
}

func getAdventures(directory string) []*adventure {
	var adventures []*adventure
	for _, f := range readDir(directory) {
		if f.IsDir() {
			checkedAdventure := checkAdventure(f.Name())
			if checkedAdventure != nil {
				adventures = append(adventures, checkedAdventure)
			}
		}
	}
	return adventures
}

func checkAdventure(name string) *adventure {
	if !exists(filepath.Join(name, "PRICE")) {
		return nil
	}

	price := readTextFile(filepath.Join(name, "PRICE"))
	description := readTextFile(filepath.Join(name, "description.html"))
	private := exists(filepath.Join(name, "PRIVATE"))
	devOnly := exists(filepath.Join(name, "DEVONLY"))
	features := parseJsonFile(filepath.Join(name, "features.json"), &adventureFeatures{})

	if devOnly && *prod {
		return nil
	}

	var puzzles []puzzle
	for _, f := range readDir(name) {
		if f.IsDir() {
			puzzles = append(puzzles, checkPuzzle(filepath.Join(name, f.Name())))
		}
	}

	return &adventure{
		name:        name,
		description: description,
		price:       price,
		private:     private,
		dev:         devOnly,
		puzzles:     puzzles,
		background:  readBinaryFile(filepath.Join(name, "background.jpg")),
		logo:        readBinaryFile(filepath.Join(name, "logo.png")),
		preview:     readBinaryFile(filepath.Join(name, "preview.png")),
		features:    features,
	}
}

func checkPuzzle(folder string) puzzle {
	parts := strings.SplitN(filepath.Base(folder), " - ", 2)
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Puzzle folder has invalid number: %s (%v)", folder, err)
	}
	name := parts[1]
	answers := readTextFileLines(filepath.Join(folder, "answers.txt"))
	var hints [][2]string
	for _, h := range readTextFileLines(filepath.Join(folder, "hints.txt")) {
		hintParts := strings.SplitN(h, ": ", 2)
		hints = append(hints, [2]string{hintParts[0], hintParts[1]})
	}
	content := replaceVariables(folder, readTextFile(filepath.Join(folder, "puzzle.html")))
	story := replaceVariables(folder, readTextFile(filepath.Join(folder, "information.html")))
	information := replaceVariables(folder, readTextFile(filepath.Join(folder, "story.html")))

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

func parseJsonFile[T interface{}](name string, parsed *T) *T {
	dataString := readTextFile(name)
	err := json.Unmarshal([]byte(dataString), parsed)
	if err != nil {
		log.Fatalf("Unable to parse features: %s", err)
	}
	return parsed
}

func replaceVariables(content, path string) (output string) {
	matcher := regexp.MustCompile(`\$(.*?)\$`)
	matches := matcher.FindAllString(content, -1)
	if matches != nil {
		for _, match := range matches {
			filename := match[1 : len(match)-1]
			if len(filename) == 0 {
				log.Fatalf("Invalid variable in file: %s: %s", match, path)
			}
			fileBytes := readBinaryFile(filepath.Join(path, filename))
			contentType := http.DetectContentType(fileBytes)
			b64 := base64.StdEncoding.EncodeToString(fileBytes)
			datauri := fmt.Sprintf("data:%s;base64,%s", contentType, b64)
			output = strings.ReplaceAll(content, match, datauri)
		}
	}
	return output
}

func readBinaryFile(name string) []byte {
	b, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("Unable to read file %s: %v", name, err)
	}
	return b
}

func readTextFile(name string) string {
	return string(readBinaryFile(name))
}

func readTextFileLines(name string) []string {
	return strings.Split(strings.TrimSpace(readTextFile(name)), "\n")
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func readDir(name string) []os.DirEntry {
	files, err := os.ReadDir(name)
	if err != nil {
		panic(err)
	}
	return files
}
