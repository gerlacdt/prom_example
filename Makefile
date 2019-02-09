PROJECT="github.com/gerlacdt/prom_example"
BINARY="server"


.PHONY=build
build: lint
	go build ${PROJECT}/cmd/server


.PHONY=lint
lint:
	golint ${PROJECT}/... && errcheck ${PROJECT}/...


.PHONY=run
run:
	go run ./cmd/server/main.go


.PHONY=clean
clean:
	rm -rf server