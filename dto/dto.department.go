package dto

type RequestDepartment struct {
	DepartmentName string `json:"name" validate:"required,min=4,max=33"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
}

type ResponseSingleDepartment struct {
	DepartmentID   string `json:"departmentId,omitempty"`
	DepartmentName string `json:"name"`
}

type ResponseMultipleDepartments struct {
	Departments []string `json:"departments"`
}
