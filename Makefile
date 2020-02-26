run-in-docker: build-linux
	docker build -t test .
	docker run --rm -it -p 80:80 -v $$PWD:/app test

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

build-ci: build-linux
	docker build -t weaming/webshot .
