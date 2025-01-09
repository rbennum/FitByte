package dto

const (
	GenderMale   = "male"
	GenderFemale = "female"

	DefaultLimit  = 5
	DefaultOffset = 0

	IdentityNumberMinLength = 5
)

type EmployeePayload struct {
	IdentityNumber   string `json:"identityNumber" validate:"required,min=5,max=33"`
	Name             string `json:"name" validate:"required,min=4,max=33"`
	EmployeeImageUri string `json:"employeeImageUri" validate:"required,uri"`
	Gender           string `json:"gender" validate:"required"`
	DepartmentID     string `json:"departmentId" validate:"required,uuid"`
}

type GetEmployeesRequest struct {
	Limit          int    `query:"limit" validate:"gte=0"`
	Offset         int    `query:"offset" validate:"gte=0"`
	IdentityNumber string `query:"identityNumber" validate:""` // validate is not set to `uuid` due to it allows wildcard
	Name           string `query:"name" validate:""`
	Gender         string `query:"gender" validate:""`
	DepartmentID   string `query:"departmentId" validate:"omitempty,uuid"`
	ManagerID      string `query:"managerId" validate:"omitempty,uuid"`
}

type UpdateEmployeeRequest struct {
	IdentityNumber   string `json:"identityNumber" validate:"omitempty,min=5,max=33"`
	Name             string `json:"name" validate:"omitempty,min=4,max=33"`
	EmployeeImageUri string `json:"employeeImageUri" validate:"omitempty,uri"`
	Gender           string `json:"gender" validate:"omitempty"`
	DepartmentID     string `json:"departmentId" validate:"omitempty,uuid"`
}
