package dialog_repo

import (
	"context"
	"database/sql"
	"github.com/GalahadKingsman/messenger_dialog/internal/models"
	"time"
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

func (s *Repo) GetUserDialogs(ctx context.Context, userID int32, limit, offset int) ([]*models.DialogInfo, error) {
	rows, err := s.db.QueryContext(ctx, `
        SELECT 
            d.id,
            u.id,
            u.login,
            (SELECT text 
             FROM messages 
             WHERE dialog_id = d.id 
             ORDER BY created_at DESC 
             LIMIT 1) as last_message,
            (SELECT created_at 
             FROM messages 
             WHERE dialog_id = d.id 
             ORDER BY created_at DESC 
             LIMIT 1) as last_activity
        FROM dialogs d
        JOIN user_dialogs_links udl ON d.id = udl.dialog_id
        JOIN users u ON u.id = udl.user_id AND u.id != $1
        WHERE udl.user_id = $1
        LIMIT $2 OFFSET $3`,
		userID, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogs []*models.DialogInfo
	for rows.Next() { // Итерируем по результатам
		var di models.DialogInfo
		var lastActivity time.Time

		// Сканируем данные из строки в структуру:
		err := rows.Scan(
			&di.ID,          // ID диалога
			&di.PeerID,      // ID собеседника
			&di.PeerLogin,   // Логин собеседника
			&di.LastMessage, // Текст последнего сообщения
			&lastActivity,   // Время последнего сообщения
		)
		if err != nil {
			return nil, err
		}

		di.LastActivity = lastActivity // Сохраняем время
		dialogs = append(dialogs, &di) // Добавляем в результат
	}

	return dialogs, nil
}
