package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/levensspel/go-gin-template/dto"
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
