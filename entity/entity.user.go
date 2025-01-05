package entity

type User struct {
	Id        string `json:"id" gorm:"primaryKey;not null;"`
	Name      string `json:"name" gorm:"not null;size:50"`
	Username  string `json:"username" gorm:"unique;not null;size:30"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime;not null"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime;not null"`
}
