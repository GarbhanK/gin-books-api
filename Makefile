
build:
	go build -o bin/main src/main/main.go

run:
	go run src/main/main.go

test:
	go test ./...

memorydb:
	go run src/main/main.go -db memorydb

postgres:
	go run src/main/main.go -db postgres

firestore:
	go run src/main/main.go -db firestore

