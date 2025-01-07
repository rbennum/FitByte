package entity

import "database/sql"

type User struct {
	Id        string         `json:"id"`
	Name      sql.NullString `json:"name"`
	Username  sql.NullString `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	UpdatedAt int64          `json:"updated_at"`
	CreatedAt int64          `json:"created_at"`
}

type UserTransactDB struct {
	Id        string         `json:"id"`
	Name      sql.NullString `json:"name"`
	Username  sql.NullString `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	UpdatedAt int64          `json:"updated_at"`
	CreatedAt int64          `json:"created_at"`
}

type GetProfile struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}
