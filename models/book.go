package models

type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type FirestoreBook struct {
	ID		int   `firestore:"id"`
	Title	string `firestore:"title"`
	Author	string `firestore:"author"`
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
    Timestamp 		string `json:"timestamp"`
    APIStatus 		string `json:"status"`
	FirestoreStatus bool   `json:"firestore_status"`
}
