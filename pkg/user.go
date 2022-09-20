package auth

type User struct {
	ID      string `gorm:"primarykey" json:"id"`
	Subject string `gorm:"unique"`
	Email   string `json:"email"`
}
