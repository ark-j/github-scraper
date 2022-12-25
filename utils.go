package githubscrape

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func CreateFilter(typef, langf, sortf, entity string, org bool) string {
	if org {
		return fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
	}
	return fmt.Sprintf("https://github.com/%s?tab=repositories&q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
}

// creates json file for per org
func CreateFile(path string, ch <-chan *Repo) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	var ll []*Repo
	for i := range ch {
		ll = append(ll, i)
	}
	b, _ := json.MarshalIndent(map[string]any{"count": len(ll), "repos": ll}, "", "  ")
	f.Write(b)
}

// cleans the string
func ClearString(s string) string {
	return strings.Trim(
		strings.ReplaceAll(
			s,
			"\n",
			"",
		),
		" ",
	)
}
