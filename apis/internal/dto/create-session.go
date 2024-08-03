package dto

type CreateSessionDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
