package repository

import "context"

const selectUser = `-- name: selectUser :one
 SELECT id, username, email FROM cas_users WHERE id = $1`

type SelectUserRows struct {
	UserID   int32  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

func (q *Queries) SelectUser(ctx context.Context, userID int32) (SelectUserRows, error) {
	row := q.db.QueryRowContext(ctx, selectUser, userID)
	var i SelectUserRows
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Email,
	)
	return i, err
}
