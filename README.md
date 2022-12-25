# scraper for github repos
- github scraper is simple github repos scraper which generate result based on orgnization or user id provided.
- it will generate json file.
- clone this repo
- use following command where `-org` is orgnizationid and `-user` is userid
- results will be generated in respective folders

```shell
cd github-scrapper
go mod tidy
go run main/main.go -org vuejs -user graydon
```
- or you can build it
```shell
go build -o github-scrapper main/main.go
```
- you can add filter
```shell
go run main/main.go -user spf13 -type source -lang go -sort stargazers
```
- check filter list in help section
```shell
go run main/main.go -help
  
  -lang string
        language filter for repositories
        possible values -> go, html, javascript, java, rust,
        python, typescript, css, haskell, shell, c++, c, ruby
  -org string
        github orgname for scraping information
  -sort string
        sort filter for repositories
        possible value -> name, stargazers
        leave empty for last updated sort
  -type string
        type filter for repositories
        possible values -> source, forks, archived, mirror, template
  -user string
        github username for scraping information
```