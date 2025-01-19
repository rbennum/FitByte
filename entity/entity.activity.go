package entity

type Activity struct {
	ActivityId        *string
	ActivityType      *string
	DoneAt            *int64
	DurationInMinutes *int64
	CaloriesBurned    *int64
	CreatedAt         *int64
}
