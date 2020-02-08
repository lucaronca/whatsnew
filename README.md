# Whatsnew
Golang lambda project to generate, maintain rss feed and sync it to S3

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
./run local <function_name>
```
Remember to create a "tmp" folder with correct file to simultate the S3 object fetching if you run the "generator" function