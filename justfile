build:
	go build -o bin/sse .

run: build
	./bin/sse
