package repository

import (
	"github.com/TimDebug/FitByte/entity"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type ActivityRepository struct {
	db *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) ActivityRepository {
	return ActivityRepository{db: db}
}

func NewActivityRepositoryInject(i do.Injector) (ActivityRepository, error) {
	db := do.MustInvoke[*pgxpool.Pool](i)
	return NewActivityRepository(db), nil
}

func (r *ActivityRepository) GetAll(
	ctx *gin.Context,
	queryArgs []interface{},
) ([]entity.Activity, error) {
	query := `
		SELECT id, activity_type, done_at, duration_in_minutes, calories_burned, created_at
		FROM activities
		WHERE
			user_id = $1
			AND ($2::TEXT IS NULL OR activity_type = $2)
			AND ($3::TIMESTAMP IS NULL OR done_at >= $3)
			AND ($4::TIMESTAMP IS NULL OR done_at <= $4)
			AND ($5::NUMERIC IS NULL OR calories_burned >= $5)
			AND ($6::NUMERIC IS NULL OR calories_burned <= $6)
		ORDER BY done_at DESC
		LIMIT $7 OFFSET $8;
	`
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []entity.Activity
	for rows.Next() {
		var activity entity.Activity
		if err := rows.Scan(
			&activity.ActivityId,
			&activity.ActivityType,
			&activity.DoneAt,
			&activity.DurationInMinutes,
			&activity.CaloriesBurned,
			&activity.CreatedAt,
		); err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
