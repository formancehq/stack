package auth

type User struct {
	ID      string `json:"id"`
	Subject string `gorm:"primarykey"`
	Email   string `json:"email"`
}
