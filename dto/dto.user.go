package dto

const (
	Create string = "create"
	Login  string = "login"
)

type UserRequestPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Action   string `json:"action" validate:"required"`
}

type RequestRegister struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UserRequestUpdate struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RequestUpdateProfile struct {
	Email           *string `json:"email" validate:"required,email"`
	Name            *string `json:"name" validate:"required,min=4,max=52"`
	UserImageUri    *string `json:"userImageUri" validate:"required,uri_with_path"`
	CompanyName     *string `json:"companyName" validate:"required,min=4,max=52"`
	CompanyImageUri *string `json:"companyImageUri" validate:"required,uri_with_path"`
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

type ResponseGetProfile struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}
