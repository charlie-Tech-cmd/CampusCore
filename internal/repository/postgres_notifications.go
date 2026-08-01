package repository

import (
	"database/sql"

	"campuscore/internal/models"
)

// PostgresNotificationRepository stores notifications in PostgreSQL.
type PostgresNotificationRepository struct {
	db *sql.DB
}

// NewPostgresNotificationRepository creates a notification repository.
func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{
		db: db,
	}
}

// Create saves a notification.
func (r *PostgresNotificationRepository) Create(
	notification *models.Notification,
) error {

	query := `
		INSERT INTO notifications
			(user_id,
			 title,
			 message,
			 type,
			 is_read)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.Type,
		notification.IsRead,
	)

	return err
}

// FindByUser returns all notifications belonging to a user.
func (r *PostgresNotificationRepository) FindByUser(
	userID string,
) ([]models.Notification, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			message,
			type,
			is_read,
			created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var notification models.Notification

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.Type,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read.
func (r *PostgresNotificationRepository) MarkAsRead(id int) error {

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}

// MarkAllAsRead marks every notification belonging to a user as read.
func (r *PostgresNotificationRepository) MarkAllAsRead(
	userID string,
) error {

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1
	`

	_, err := r.db.Exec(query, userID)
	return err
}

// Delete removes a notification.
func (r *PostgresNotificationRepository) Delete(id int) error {

	query := `
		DELETE FROM notifications
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}
