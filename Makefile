
build:
	PRIMARY_DB=memorydb ENABLE_CACHE=true \
	go build -o bin/main main/main.go

run:
	PRIMARY_DB=memorydb ENABLE_CACHE=false \
	go run main/main.go

run-multi:
	PRIMARY_DB=memorydb SECONDARY_DB=postgres ENABLE_CACHE=false \
	go run main/main.go

pg:
	podman run \
	-e POSTGRES_USER=gin \
	-e POSTGRES_PASSWORD=ginpass \
	-p 5432:5432 \
	postgres:16-alpine

test:
	go test ./...

memorydb:
	PRIMARY_DB=memorydb \
	go run main/main.go -db memorydb

postgres:
	PRIMARY_DB=postgres \
	go run main/main.go -db postgres

firestore:
	PRIMARY_DB=firestore \
	go run main/main.go -db firestore
