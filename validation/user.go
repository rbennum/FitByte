package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/levensspel/go-gin-template/dto"
	repository "github.com/levensspel/go-gin-template/repository/user"
)

var validate = func() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("uri_with_path", IsValidURI)

	return v
}()

func ValidateUserCreate(input dto.UserRequestPayload, r repository.UserRepository) error {
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	// (!) kedua dibawah ini transaksi db
	// err = ValidateIsUsernameExist(input.Username, r)
	// if err != nil {
	// 	return err
	// }
	// err = ValidateIsEmailExist(input.Email, r)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func ValidateUserLogin(input dto.UserRequestPayload) error {
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}

	return nil
}

func ValidateUpdateProfile(input dto.RequestUpdateProfile) error {
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	return nil
}
