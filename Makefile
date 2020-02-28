.PHONY: build clean deploy gomodgen run-generator-local

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/generator cmd/generator/main.go cmd/generator/local.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/uploader cmd/uploader/main.go cmd/uploader/local.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

run-generator-local:
	env LOCAL=true go run cmd/generator/main.go cmd/generator/local.go
