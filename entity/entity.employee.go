package entity

type Employee struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	IdentityNumber   string `json:"identityNumber"`
	EmployeeImageUri string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentId     string `json:"departmentId"`
	CreatedAt        int64  `json:"createdAt"`
	UpdatedAt        int64  `json:"updatedAt"`
}
