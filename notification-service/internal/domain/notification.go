package domain

type NotificationType string

const (
	EmailNotification NotificationType = "email"
	SMSNotification   NotificationType = "sms"
)

type Notification struct {
	UserID  int64
	To      string
	Type    NotificationType
	Subject string
	Message string
}

type NotificationSender interface {
	Send(notification Notification) error
}

type NotificationUseCase interface {
	Send(notification Notification) error
}
