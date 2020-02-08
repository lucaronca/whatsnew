.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/generator generator/main.go generator/local.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/uploader uploader/main.go uploader/local.go uploader/fileuploader.go uploader/addfiletos3.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
