package importer

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
