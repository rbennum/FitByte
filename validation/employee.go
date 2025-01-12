package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
)

func ValidateEmployeeCreate(input *dto.EmployeePayload) error {
	if input.IdentityNumber == "" {
		return errors.New("Invalid identity number param")
	}
	if input.Name == "" {
		return errors.New("Invalid name param")
	}
	if input.EmployeeImageUri == "" {
		return errors.New("Invalid image URI param")
	}
	if input.DepartmentID == "" {
		return errors.New("Invalid department id param")
	}

	extension := input.EmployeeImageUri[len(input.EmployeeImageUri)-5:]
	if extension != ".jpeg" {
		extension = extension[1:]
		if extension != ".jpg" && extension != ".png" {
			return errors.New("Invalid image URI param")
		}
	}

	gender := strings.ToLower(input.Gender)
	// fmt.Println("gender:", strings.ToLower(gender))
	switch gender {
	case "":
		input.Gender = ""
	case dto.GenderMale:
		input.Gender = dto.GenderMale
		break
	case dto.GenderFemale:
		input.Gender = dto.GenderFemale
		break
	default:
		return errors.New("Invalid gender param")
	}

	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	return nil
}

func ValidateEmployeeGet(input *dto.GetEmployeesRequest) error {
	gender := strings.ToLower(input.Gender)
	// fmt.Println("gender:", strings.ToLower(gender))
	switch gender {
	case "":
		input.Gender = ""
	case dto.GenderMale:
		input.Gender = dto.GenderMale
		break
	case dto.GenderFemale:
		input.Gender = dto.GenderFemale
		break
	default:
		return errors.New("Invalid gender param")
	}

	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	return nil
}

func ValidateEmployeeUpdate(input *dto.UpdateEmployeeRequest) error {
	if input.IdentityNumber == "" &&
		input.Name == "" &&
		input.EmployeeImageUri == "" &&
		input.Gender == "" &&
		input.DepartmentID == "" {
		return helper.ErrBadRequest
	}

	gender := strings.ToLower(input.Gender)
	// fmt.Println("gender:", strings.ToLower(gender))
	switch gender {
	case "":
		input.Gender = ""
	case dto.GenderMale:
		input.Gender = dto.GenderMale
		break
	case dto.GenderFemale:
		input.Gender = dto.GenderFemale
		break
	default:
		return errors.New("Invalid gender param")
	}

	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	return nil
}
