
build:
	DB_TYPE=memorydb ENABLE_CACHE=true \
	go build -o bin/main main/main.go

run:
	DB_TYPE=memorydb ENABLE_CACHE=false \
	go run main/main.go

test:
	go test ./...

memorydb:
	DB_TYPE=memorydb \
	go run main/main.go -db memorydb

postgres:
	DB_TYPE=postgres \
	go run main/main.go -db postgres

firestore:
	DB_TYPE=firestore \
	go run main/main.go -db firestore
