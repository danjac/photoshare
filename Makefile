GOPATH := ${PWD}:${GOPATH}
export GOPATH

#build: build-go migrate build-ui
build: build-go build-ui

build-go:
	#godep restore
	go build -o bin/serve -i commands/server/main.go
	go build -o bin/import -i commands/import/main.go


build-ui: 
	npm install
	bower install
	gulp build

#migrate:
#	go get bitbucket.org/liamstask/goose/cmd/goose
#	goose up

test:
	go test ./...

