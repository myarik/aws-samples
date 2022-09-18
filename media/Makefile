.PHONY: build clean deploy

functions := $(shell find functions -name \*main.go | awk -F'/' '{print $$2}')

build: clean
	export GO111MODULE=on
	@for function in $(functions) ; do \
		env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$$function functions/$$function/main.go ; \
    done

clean:
	rm -rf ./bin

deploy:
	sls deploy --verbose

deploy-function:
	sls deploy function --function $(filter-out $@,$(MAKECMDGOALS))

