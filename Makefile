BINARY_NAME=meal-planner

build:
	go mod tidy
	rm -rf dist || true
	mkdir dist
	go build -o dist/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf dist || true

.PHONY: all test clean
