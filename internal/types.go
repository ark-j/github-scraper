package internal

type Repo struct {
	Title       string
	Link        string
	Description string
	Language    string
	Forks       string
	Stars       string
}

type Filter struct {
	Type string
	Lang string
	Sort string
}
