package internal

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const URL = "https://github.com"

// creates json file for per org
func CreateFile(path string, ch <-chan *Repo) {
	f, err := os.Create(path)
	if err != nil {
		log.Println("ERROR", "msg=can't create file", "error=", err)
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
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}
