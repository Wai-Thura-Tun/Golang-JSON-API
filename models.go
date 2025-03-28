package main

type Status string

const (
	Read    Status = "read"
	Reading Status = "reading"
	ToRead  Status = "to_read"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"_"`
}

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Status Status `json:"status" gorm:"default:to_read"`
	Author string `json:"author"`
	Year   uint   `json:"year"`
	UserID uint   `json:"user_id"`
}
