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
## TODO
- add filter