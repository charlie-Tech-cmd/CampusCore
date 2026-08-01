package services

import (
	"campuscore/internal/models"
	"campuscore/internal/notification"
	"fmt"
)

type NotificationService struct {
	repo   models.NotificationRepository
	worker *notification.Worker
}

func NewNotificationService(
	repo models.NotificationRepository,
	worker *notification.Worker,
) *NotificationService {

	return &NotificationService{
		repo:   repo,
		worker: worker,
	}
}

// CreateNotification creates and dispatches a notification.
func (s *NotificationService) CreateNotification(
	userID string,
	title string,
	message string,
	notificationType models.NotificationType,
) error {

	n := &models.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notificationType,
		IsRead:  false,
	}

	if err := s.repo.Create(n); err != nil {
		return err
	}

	if s.worker != nil {
		s.worker.Enqueue(notification.Task{
			RecipientID: userID,
			Type:        string(notificationType),
			Payload:     message,
		})
	}

	return nil
}

// NotifyPayment sends a payment confirmation notification.
func (s *NotificationService) NotifyPayment(
	userID string,
	amount float64,
) error {

	title := "Payment Successful"

	message := fmt.Sprintf(
		"Your payment of ₦%.2f has been received successfully.",
		amount,
	)

	return s.CreateNotification(
		userID,
		title,
		message,
		models.NotificationPayment,
	)
}

// NotifyClearance sends a clearance update.
func (s *NotificationService) NotifyClearance(
	userID string,
	office string,
) error {

	title := "Clearance Updated"

	message := fmt.Sprintf(
		"%s has approved your clearance.",
		office,
	)

	return s.CreateNotification(
		userID,
		title,
		message,
		models.NotificationClearance,
	)
}

// NotifyTranscriptReady informs a student that a transcript is ready.
func (s *NotificationService) NotifyTranscriptReady(
	userID string,
) error {

	return s.CreateNotification(
		userID,
		"Transcript Ready",
		"Your transcript has been generated successfully.",
		models.NotificationTranscript,
	)
}
