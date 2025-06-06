package dialogservice

import (
	"context"
	"database/sql"
)

func (s *sql.DB) CreateDialog(ctx context.Context, userID, peerID int32) (int32, error) {
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
