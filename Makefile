PROJECT="github.com/gerlacdt/prom_example"

build: lint
	go build ${PROJECT}/...


lint:
	golint ${PROJECT}/... && errcheck ${PROJECT}/...
