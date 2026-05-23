package repository

import (
	"context"

	"github.com/google/uuid"
)

const getActivityPillars = `
SELECT pillar FROM activity_pillars WHERE activity_id = $1
`

func (q *Queries) GetActivityPillars(ctx context.Context, activityID uuid.UUID) ([]Pillar, error) {
	rows, err := q.db.Query(ctx, getActivityPillars, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pillars []Pillar
	for rows.Next() {
		var p Pillar
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		pillars = append(pillars, p)
	}
	return pillars, rows.Err()
}

const deleteActivityPillars = `
DELETE FROM activity_pillars WHERE activity_id = $1
`

func (q *Queries) DeleteActivityPillars(ctx context.Context, activityID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteActivityPillars, activityID)
	return err
}
