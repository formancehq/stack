package auth

type User struct {
	Subject string `gorm:"primarykey"`
	Email   string
}
