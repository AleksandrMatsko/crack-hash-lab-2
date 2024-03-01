.PHONY: build run test vet

# build project
build:
	go build -v distributed.systems.labs/manager

# build and run project
run: build
	./manager

# run tests recursively with data race detection
test:
	go test -race ./...

# vet project
vet:
	go vet ./...
