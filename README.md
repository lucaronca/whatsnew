# Whatsnew
Golang lambda project to generate, maintain rss feed and sync it to S3 <br>
Synched Rss feed can then be used to automate push notifications throught bot services like ifttt or automate.io

## Project structure
![Project_structure](https://whatsnew-bucket.s3.us-east-1.amazonaws.com/structure.png)

## build
```bash
make build
```
## deploy
```bash
sls deploy
```
## run locally
```bash
./run_local <function_name>
```
Remember to create a "tmp" folder with correct file to simultate the S3 object fetching if you run the "generator" function
## run tests
```bash
go test ./...
```
