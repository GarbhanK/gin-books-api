
build:
	go build -o bin/main main/main.go

run:
	go run main/main.go

memorydb:
	go run main/main.go -db memorydb

postgres:
	go run main/main.go -db postgres

firestore:
	go run main/main.go -db firestore

