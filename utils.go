package githubscrape

import "fmt"

func CreateFilter(typef, langf, sortf, entity string, org bool) string {
	if org {
		return fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
	}
	return fmt.Sprintf("https://github.com/%s?tab=repositories&q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
}
