package usecase

import (
	"notification-service/internal/domain"
)

type NotificationUseCase struct {
	sender domain.NotificationSender
}

func NewNotificationUseCase(sender domain.NotificationSender) *NotificationUseCase {
	return &NotificationUseCase{sender: sender}
}

func (uc *NotificationUseCase) Send(notification domain.Notification) error {
	return uc.sender.Send(notification)
}
