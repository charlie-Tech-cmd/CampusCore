package models

import "time"

// NotificationType represents the category of notification.
type NotificationType string

const (
	NotificationPayment      NotificationType = "payment"
	NotificationClearance    NotificationType = "clearance"
	NotificationTranscript   NotificationType = "transcript"
	NotificationResult       NotificationType = "result"
	NotificationAnnouncement NotificationType = "announcement"
	NotificationAdmission    NotificationType = "admission"
)

// Notification represents a user notification stored in the system.
type Notification struct {
	ID        int              `json:"id" db:"id"`
	UserID    string           `json:"user_id" db:"user_id"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Type      NotificationType `json:"type" db:"type"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

// NotificationRepository defines notification persistence.
type NotificationRepository interface {
	Create(notification *Notification) error
	FindByUser(userID string) ([]Notification, error)
	MarkAsRead(id int) error
	MarkAllAsRead(userID string) error
	Delete(id int) error
}
