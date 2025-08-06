package sender

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/domain"
	"github.com/Wrestler094/scalable-ecommerce-platform/notification-service/internal/infrastructure/sender/dto"
)

type EmailSender struct {
	apiKey    string
	fromEmail string
	fromName  string
}

func NewEmailSender(apiKey, fromEmail, fromName string) *EmailSender {
	return &EmailSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *EmailSender) Send(n domain.Notification) error {
	const op = "emailSender.Send"

	if n.Type != domain.EmailNotification {
		return fmt.Errorf("%s: unsupported notification type", op)
	}

	form := url.Values{}
	form.Set("apikey", s.apiKey)
	form.Set("from", s.fromEmail)
	form.Set("fromName", s.fromName)
	form.Set("to", n.To)
	form.Set("subject", n.Subject)
	form.Set("bodyText", n.Message)
	form.Set("isTransactional", "true")

	req, err := http.NewRequest("POST", "https://api.elasticemail.com/v2/email/send", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("%s: failed to create request: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: failed to send request: %w", op, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: failed with status %d: %s", op, resp.StatusCode, string(body))
	}

	var result dto.ElasticResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%s: failed to decode response: %w", op, err)
	}

	if !result.Success {
		return fmt.Errorf("%s: failed to send via Elastic Email: %s", op, result.Error)
	}

	return nil
}
