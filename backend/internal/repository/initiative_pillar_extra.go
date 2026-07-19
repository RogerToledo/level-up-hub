package repository

import (
	"context"

	"github.com/google/uuid"
)

const getInitiativePillars = `
SELECT pillar FROM initiative_pillars WHERE initiative_id = $1
`

func (q *Queries) GetInitiativePillars(ctx context.Context, initiativeID uuid.UUID) ([]Pillar, error) {
	rows, err := q.db.Query(ctx, getInitiativePillars, initiativeID)
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

const deleteInitiativePillars = `
DELETE FROM initiative_pillars WHERE initiative_id = $1
`

func (q *Queries) DeleteInitiativePillars(ctx context.Context, initiativeID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteInitiativePillars, initiativeID)
	return err
}
