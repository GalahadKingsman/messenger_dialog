package dialog_repo

import (
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (s *Repo) CreateDialog(ctx context.Context, userID, peerID int32) (int32, error) {
	var dialogID int32
	err := s.db.QueryRowContext(ctx, `
		WITH new_dialog AS (
			INSERT INTO dialogs DEFAULT VALUES RETURNING id
		),
		insert_links AS (
			INSERT INTO user_dialogs_links (user_id, dialog_id)
			SELECT $1, id FROM new_dialog
			UNION ALL
			SELECT $2, id FROM new_dialog
			RETURNING dialog_id
		)
		SELECT dialog_id FROM insert_links LIMIT 1`,
		userID, peerID).Scan(&dialogID)

	return dialogID, err
}

func (s *Repo) GetUserDialogs(ctx context.Context, userID int32, limit, offset int) ([]*DialogInfo, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT d.id, u.id, u.login, 
			(SELECT text FROM messages WHERE dialog_id = d.id ORDER BY created_at DESC LIMIT 1),
			(SELECT created_at FROM messages WHERE dialog_id = d.id ORDER BY created_at DESC LIMIT 1)
		FROM dialogs d
		JOIN user_dialogs_links udl ON d.id = udl.dialog_id
		JOIN users u ON u.id = udl.user_id AND u.id != $1
		WHERE udl.user_id = $1
		LIMIT $2 OFFSET $3`, userID, limit, offset)

	return dialogs, nil
}
