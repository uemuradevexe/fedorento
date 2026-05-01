package data

// Topic represents a documentation leaf node with content.
type Topic struct {
	ID          string
	Title       string
	Description string
	Code        string
	Explanation string
	Language    string // "php" for most Laravel content
}

// Section groups related topics under a nav node.
type Section struct {
	ID     string
	Title  string
	Topics []Topic
}

// Chapter is a top-level entry (e.g. "Laravel 13").
type Chapter struct {
	ID       string
	Title    string
	Sections []Section
}
