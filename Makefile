build:
	go build

test:
	go test -v ./... --cover

lint:
	go fmt	./...

coverage:
	go test -coverprofile=coverage.txt ./...

html: coverage
	go tool cover -html=coverage.txt