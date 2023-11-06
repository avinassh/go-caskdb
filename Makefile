build:
	go build

test:
	go test -v ./... --cover

bench:
	go test -bench=. -run=^# -benchtime=20x 

lint:
	go fmt	./...

coverage:
	go test -coverprofile=coverage.txt ./...

html: coverage
	go tool cover -html=coverage.txt