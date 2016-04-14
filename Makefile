test-fast:
	go test ./...
.PHONY: test-fast

test:
	go test -race ./...
.PHONY: test

bench:
	go test -bench=. -benchmem ./...
.PHONY: bench

install:
	go install ./...
.PHONY: install


tools:
	go get -u github.com/onsi/ginkgo/...
	go get -u github.com/tools/godep
	go get -u github.com/smartystreets/goconvey
.PHONY: tools