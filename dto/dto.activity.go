package dto

type ResponseActivity struct {
	Id                string `json:"activityId"`
	ActivityType      string `json:"activityType"`
	DoneAt            string `json:"doneAt"`
	DurationInMinutes int    `json:"durationInMinutes"`
	CaloriesBurned    int    `json:"caloriesBurned"`
	CreatedAt         string `json:"createdAt"`
}
