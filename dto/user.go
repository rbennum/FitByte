package dto

type RequestRegister struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserRequestUpdate struct {
    Id    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Responses
type Response struct {
	Id       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}
