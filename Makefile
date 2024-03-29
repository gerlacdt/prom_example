PROJECT="github.com/gerlacdt/prom_example"
BINARY="server"


.PHONY=build
build: server-build client-build

.PHONY=server-build
server-build: lint
	go build ${PROJECT}/cmd/server


.PHONY=client-build
client-build:
	go build ${PROJECT}/cmd/client


.PHONY=lint
lint:
	golint ${PROJECT}/... && errcheck ${PROJECT}/...


.PHONY=run
run: build
	go run ./cmd/server/main.go

.PHONY=test
test:
	go test ${PROJECT}/...


.PHONY=clean
clean:
	rm -rf server
