package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
)

func ValidateEmployeeCreate(input *dto.EmployeePayload) error {
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
