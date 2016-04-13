test-fast:
	go test ./...
.PHONY: test-fast

test:
	go test -race ./...
.PHONY: test

bench:
	go test -bench=. -benchmem ./...
.PHONY: bench

tools:
	go get -u github.com/onsi/ginkgo/...
	go get -u github.com/tools/godep