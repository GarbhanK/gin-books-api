package models

type Book struct {
	Id	   string `json:"id" firestore:"id"`
	Title  string `json:"title" firestore:"title"`
	Author string `json:"author" firestore:"author"`
}

type Status struct {
	Timestamp string `json:"timestamp"`
	APIStatus string `json:"api_status"`
	DBStatus  string `json:"database_status"`
	DBType    string `json:"database_type"`
}

type InsertBookInput struct {
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

