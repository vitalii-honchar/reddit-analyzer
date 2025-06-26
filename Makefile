
build:
	go build -o out/redditanalyzer ./cmd/redditanalyzer

lint:
	go vet ./...
