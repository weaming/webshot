run-in-docker:
	docker build -t test .
	docker run --rm -it -p 80:80 -v $$PWD:/app test

build-ci:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
	docker build -t weaming/webshot .
