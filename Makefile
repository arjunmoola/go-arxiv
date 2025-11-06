install: build
	go install github.com/arjunmoola/go-arxiv/cmd/garx

build:
	go build -o bin/ github.com/arjunmoola/go-arxiv/cmd/garx

clean:
	go clean -i github.com/arjunmoola/go-arxiv/cmd/garx
	rm -rf bin/
