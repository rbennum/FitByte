package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/levensspel/go-gin-template/dto"
	repository "github.com/levensspel/go-gin-template/repository/user"
)

var validate = validator.New()

func ValidateIsEmailExist(email string, r repository.UserRepository) error {
	err := r.IsEmailRegistered(email, nil)
	return err
}

func ValidateIsUsernameExist(username string, r repository.UserRepository) error {
	err := r.IsUsernameRegistered(username, nil)
	return err
}

func ValidateUserCreate(input dto.RequestRegister, r repository.UserRepository) error {
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	err = ValidateIsUsernameExist(input.Username, r)
	if err != nil {
		return err
	}
	err = ValidateIsEmailExist(input.Email, r)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUserLogin(input dto.RequestLogin) error {
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}

	return nil
}
