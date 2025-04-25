# gin-books-api

- A REST API using the `Gin` web framework for managing a collection of books with modular backends
    - Firestore is a cloud SDK backends
    - MemoryDB maps the data to an in-memory store, useful for testing
    - Postgres maps the data to a generic SQL backend
- Makes use of GCPs generous free tier
- Uses [Google Firestore](https://cloud.google.com/firestore?hl=en) for a scalable document database
- [Started from this article](https://blog.logrocket.com/rest-api-golang-gin-gorm/)

## TODOs
- get Postgres interface working
- see if I can/need to write tests for the Firestore/Postgres interfaces
- Have it so multiple db endpoints can be selected
    - have the same data be inserted to Postgres & MemoryDB at the same time
    - use channels with separate queues to split messages between them
