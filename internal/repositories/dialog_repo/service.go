package dialog_repo

import (
	"context"
	"database/sql"
	"fmt"
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

// Создает новый диалог в базе данных
func (s *Repo) CreateDialog(userID, peerID int32, dialogName string) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Вставляем новый диалог
	var dialogID int
	err = tx.QueryRow(`
        INSERT INTO dialogs (name) VALUES ($1) RETURNING id
    `, dialogName).Scan(&dialogID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert dialog: %v", err)
	}

	// Связываем пользователей с диалогом
	_, err = tx.Exec(`
        INSERT INTO users_dialogs_links (user_id, dialog_id, link_name) 
        VALUES ($1, $2, $3), ($4, $2, $5)
    `, userID, dialogID, dialogName, peerID, dialogName)
	if err != nil {
		return 0, fmt.Errorf("failed to link users to dialog: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return dialogID, nil
}

func (s *Repo) CheckDialog(userID, peerID int32) (int, string, error) {
	query := `
		SELECT d.id, d.name 
		FROM dialogs d
		JOIN users_dialogs_links ud1 ON d.id = ud1.dialog_id
		JOIN users_dialogs_links ud2 ON d.id = ud2.dialog_id
		WHERE ud1.user_id = $1 AND ud2.user_id = $2
		LIMIT 1
	`

	var dialogID int
	var dialogName string

	err := s.db.QueryRow(query, userID, peerID).Scan(&dialogID, &dialogName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil // диалог не найден - это не ошибка
		}
		return 0, "", fmt.Errorf("failed to query existing dialog: %v", err)
	}

	return dialogID, dialogName, nil
}

func (r *Repo) GetPeerID(dialogID, userID int32) (int32, error) {
	const q = `
        SELECT user_id
          FROM users_dialogs_links
         WHERE dialog_id = $1
           AND user_id <> $2
         LIMIT 1
    `
	var peerID int32
	err := r.db.QueryRow(q, dialogID, userID).Scan(&peerID)
	if err != nil {
		return 0, fmt.Errorf("GetPeerID: %w", err)
	}
	return peerID, nil
}

func (r *Repo) GetUserDialogs(userID int32, limit, offset int32) ([]*models.DialogInfo, error) {
	query := `
		SELECT 
			d.id AS dialog_id,
			CASE WHEN ud1.user_id = $1 THEN ud2.user_id ELSE ud1.user_id END AS peer_id,
			u.login AS peer_login,
			m.text AS last_message
		FROM dialogs d
		JOIN users_dialogs_links ud1 ON d.id = ud1.dialog_id
		JOIN users_dialogs_links ud2 ON d.id = ud2.dialog_id AND ud2.user_id != $1
		JOIN users u ON u.id = CASE WHEN ud1.user_id = $1 THEN ud2.user_id ELSE ud1.user_id END
		LEFT JOIN LATERAL (
			SELECT text 
			FROM messages 
			WHERE dialog_id = d.id 
			ORDER BY create_date DESC 
			LIMIT 1
		) m ON true
		WHERE ud1.user_id = $1
		ORDER BY d.create_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query dialogs: %w", err)
	}
	defer rows.Close()

	var dialogs []*models.DialogInfo
	for rows.Next() {
		var (
			dialogID    int32
			peerID      int32
			peerLogin   string
			lastMessage sql.NullString
		)

		if err := rows.Scan(&dialogID, &peerID, &peerLogin, &lastMessage); err != nil {
			return nil, fmt.Errorf("failed to scan dialog row: %w", err)
		}

		dialog := &models.DialogInfo{
			ID:          dialogID,
			PeerID:      peerID,
			PeerLogin:   peerLogin,
			LastMessage: "",
		}

		if lastMessage.Valid {
			dialog.LastMessage = lastMessage.String
		}

		dialogs = append(dialogs, dialog)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return dialogs, nil
}

func (r *Repo) SendMessage(dialogID, userID int32, text string) (int32, time.Time, error) {
	var messageID int32
	var createdAt time.Time

	query := `
		INSERT INTO messages (dialog_id, user_id, text, create_date)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, create_date
	`

	err := r.db.QueryRow(query, dialogID, userID, text).Scan(&messageID, &createdAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to insert message: %w", err)
	}

	return messageID, createdAt, nil
}

func (r *Repo) GetDialogMessages(ctx context.Context, dialogID int32, limit, offset *int32) ([]*models.Message, error) {
	query := `
        SELECT id, user_id, text, create_date 
        FROM messages 
        WHERE dialog_id = $1 
        ORDER BY create_date DESC`

	args := []interface{}{dialogID}
	paramIndex := 2

	if limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", paramIndex)
		args = append(args, *limit)
		paramIndex++
	}
	if offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", paramIndex)
		args = append(args, *offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Text, &msg.CreateDate); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return messages, nil
}
