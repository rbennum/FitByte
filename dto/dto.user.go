package dto

type UserRequestPayload struct {
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
type ResponseAuth struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResponseGetProfile struct {
	Preference *string `json:"preference"`
	WeightUnit *string `json:"weightUnit"`
	HeightUnit *string `json:"heightUnit"`
	Weight     *int    `json:"weight"`
	Height     *int    `json:"height"`
	Email      string  `json:"email"`
	Name       *string `json:"name"`
	ImageUri   *string `json:"imageUri"`
}
