package dto

const (
	Create string = "create"
	Login  string = "login"
)

type UserRequestPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Action   string `json:"action" validate:"required"`
}

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
	Id       string `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResponseRegister struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResponseLogin struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
