package models

type Book struct {
	Title  string `json:"title" firestore:"title"`
	Author string `json:"author" firestore:"author"`
}

type Status struct {
	Timestamp       string `json:"timestamp"`
	APIStatus       string `json:"status"`
	FirestoreStatus string `json:"firestore_status"`
}

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type FindAuthorInput struct {
	Author string `json:"author"`
}
