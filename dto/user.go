package dto

type RequestRegister struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Response struct {
	Id       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}
