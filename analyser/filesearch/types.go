package filesearch

import "os"

type Match struct {
	File       string `json:"file"`
	Line       string `json:"line"`
	LineNumber uint32 `json:"line_number"`
	Match      string `json:"match"`
}

type Result struct {
	File       string   `json:"file"`
	Matches    []Match  `json:"matches"`
	ImportedBy []string `json:"imported_by"`
}

type Results map[string]Result

type FileStructure struct {
	Variables map[string][]Match `json:"variables"`
	Functions map[string][]Match `json:"functions"`
	Mixins    map[string][]Match `json:"mixins"`
}

type FindFunc func(path string) []Match

type FileData struct {
	file os.DirEntry
	path string
}

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type D3Data struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
