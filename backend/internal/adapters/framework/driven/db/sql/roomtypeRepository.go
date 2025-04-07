package sql

import (
	"context"
	"database/sql"

	"github.com/sql-project-backend/internal/models/dto"
)

type PostgresRoomTypeRepository struct {
	db *sql.DB
}

func NewPostgresRoomTypeRepository(db *sql.DB) *PostgresRoomTypeRepository {
	return &PostgresRoomTypeRepository{db: db}
}

func (r *PostgresRoomTypeRepository) ListRoomTypes(ctx context.Context) ([]*dto.RoomTypePublic, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name
        FROM room_type
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*dto.RoomTypePublic
	for rows.Next() {
		var rt dto.RoomTypePublic
		if err := rows.Scan(&rt.RoomTypeID, &rt.Name); err != nil {
			return nil, err
		}
		out = append(out, &rt)
	}
	return out, rows.Err()
}
