package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
