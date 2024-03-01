.PHONY: build run test vet

# build project
build:
	go build -v distributed.systems.labs/worker

# build and run project
run: build
	./worker

# run tests recursively with data race detection
test:
	go test -race ./...

# vet project
vet:
	go vet ./...
