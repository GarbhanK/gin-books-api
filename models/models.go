package models

type Book struct {
	Id     string `json:"id" firestore:"id"`
	Title  string `json:"title" firestore:"title"`
	Author string `json:"author" firestore:"author"`
}

type APIStatus struct {
	Timestamp string     `json:"timestamp"`
	APIStatus string     `json:"api_status"`
	DBStatus  []DBStatus `json:"db_status"`
}

type DBStatus struct {
	Tier       string `json:"tier"`
	Type       string `json:"type"`
	Connection string `json:"status"`
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
