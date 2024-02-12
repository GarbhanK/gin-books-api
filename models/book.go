package models

type Book struct {
	ID     int    `json:"id" gorm:"primary_key" firestore:"id"`
	Title  string `json:"title" firestore:"title"`
	Author string `json:"author" firestore: "author"`
}

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Status struct {
	Timestamp       string `json:"timestamp"`
	APIStatus       string `json:"status"`
	FirestoreStatus string `json:"firestore_status"`
}
