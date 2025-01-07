package dto

const (
	GenderMale   = "male"
	GenderFemale = "female"

	DefaultLimit  = 5
	DefaultOffset = 0
)

type GetEmployeesRequest struct {
	Limit          int    `query:"limit" validate:"gte=0"`
	Offset         int    `query:"offset" validate:"gte=0"`
	IdentityNumber string `query:"identityNumber" validate:""` // validate is not set to `uuid` due to it allows wildcard
	Name           string `query:"name" validate:""`
	Gender         string `query:"gender" validate:""`
	DepartmentID   string `query:"departmentId" validate:"omitempty,uuid"`
	ManagerID      string `query:"managerId" validate:"omitempty,uuid"`
}

type GetEmployeeResponseItem struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	EmployeeImageUri string `json:"employeeImageUri"`
	DepartmentID     string `json:"departmentId"`
}
