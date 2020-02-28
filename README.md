# Whatsnew
Golang lambda project to generate, maintain rss feed and sync it to S3 <br>
Synched Rss feed can then be used to automate push notifications throught bot services like ifttt or automate.io

## Project structure
![Project_structure](https://github.com/lucaronca/whatsnew/blob/master/assets/project_structure.png)

## build
```bash
make build
```
## deploy
```bash
sls deploy
```
## run locally
`Generator` function can be ran locally:
```bash
make run-generator-local
```
Remember to create a "tmp" folder with correct file to simultate the S3 object fetching if you run the "generator" function
## run tests
```bash
go test ./...
```
