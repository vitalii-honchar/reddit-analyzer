
build:
	go build -o out/redditanalyzer ./cmd/redditanalyzer

lint:
	golangci-lint run
