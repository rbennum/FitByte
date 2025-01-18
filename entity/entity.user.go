package entity

type User struct {
	Id           *string `json:"id"`
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
	Preference   *string `json:"preference"`
	WeightUnit   *string `json:"weight_unit"`
	HeightUnit   *string `json:"height_unit"`
	Weight       *int    `json:"weight"`
	Height       *int    `json:"height"`
	Name         *string `json:"name"`
	ImageUri     *string `json:"image_uri"`
	UpdatedAt    int64   `json:"updated_at"`
	CreatedAt    int64   `json:"created_at"`
}
