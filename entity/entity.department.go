package entity

type Department struct {
	Id        string `json:"department_id" gorm:"primaryKey; not null;"`
	Name      string `json:"department_name" gorm:"unique; not null; size:33"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime;not null"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime;not null"`
}
