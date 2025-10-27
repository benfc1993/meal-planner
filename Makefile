BINARY_NAME=meals

build:
	go mod tidy
	rm -rf dist || true
	mkdir dist
	cp -r views/ dist/views/
	cp -r static/ dist/static/
	go build -o dist/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf dist || true

.PHONY: all test clean
