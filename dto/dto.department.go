package dto

type RequestDepartment struct {
	DepartmentName string `json:"department_name" validate:"required,min=4,max=33"`
}

type ResponseSingleDepartment struct {
	DepartmentID   string `json:"department_id,omitempty"`
	DepartmentName string `json:"department_name"`
}

type ResponseMultipleDepartments struct {
	Departments []string `json:"departments"`
}
