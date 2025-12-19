# gin-books-api

- A REST API using the `Gin` web framework for managing a collection of books with modular backends
    - Firestore is a cloud SDK backends
    - MemoryDB maps the data to an in-memory store, useful for testing
    - Postgres maps the data to a generic SQL backend
- Makes use of GCPs generous free tier
- Uses [Google Firestore](https://cloud.google.com/firestore?hl=en) for a scalable document database
- [Started from this article](https://blog.logrocket.com/rest-api-golang-gin-gorm/)

## TODOs
- [x] get Postgres interface working
- [ ] add an `insert_timestamp` column to the schema
- [x] get memoryDB working with lowercase table names
- [ ] see if I can/need to write tests for the Firestore/Postgres interfaces
    - mock the sdk interfaces?
- [ ] see if I should move database tests into the database module
- [x] Have it so multiple db endpoints can be selected
    - have the same data be inserted to Postgres & MemoryDB at the same time
    - use channels with separate queues to split messages between them
- [ ] Integration testing with go test containers: https://golang.testcontainers.org/quickstart/

This project demonstrates both idiomatic Go CRUD API patterns and advanced Go concurrency. Writing to multiple DBs concurrently using Go channels and goroutines is deliberately implemented for demonstration, with error handling and discussion of consistency tradeoffs clearly documented. In a production system, transactional guarantees would be handled via distributed transactions or compensating logic; here, the focus is on illustrating Go's concurrency primitives within an API service
