package entity

type Department struct {
	Id        string `json:"department_id"`
	Name      string `json:"department_name"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
}
