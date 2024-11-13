package service

import (
	"bytes"
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"errors"
	"html/template"
	"time"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
	templateRepository     domain.TemplateRepository
	hub                    *dto.Hub
}

func NewNotification(notificationRepository domain.NotificationRepository, templateRepository domain.TemplateRepository, hub *dto.Hub) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
		templateRepository:     templateRepository,
		hub:                    hub,
	}
}

// FindByUser implements domain.NotificationService.
func (n *notificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	notifications, err := n.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		return nil, err
	}
	var result []dto.NotificationData
	for _, notification := range notifications {
		result = append(result, dto.NotificationData{
			ID:        notification.ID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		})
	}

	if result == nil {
		result = make([]dto.NotificationData, 0)
	}

	return result, nil
}

// Insert implements domain.NotificationService.
func (n *notificationService) Insert(ctx context.Context, userId int64, code string, data map[string]string) error {
	templ, err := n.templateRepository.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if templ == (domain.Template{}) {
		return errors.New("template not found")
	}

	body := new(bytes.Buffer)
	t := template.Must(template.New("notif").Parse(templ.Body))
	err = t.Execute(body, data)

	if err != nil {
		return err
	}

	notification := domain.Notification{
		UserID:    userId,
		Title:     templ.Title,
		Body:      body.String(),
		Status:    1,
		IsRead:    0,
		CreatedAt: time.Now(),
	}

	err = n.notificationRepository.Insert(ctx, &notification)
	if err != nil {
		return err
	}

	if channel, ok := n.hub.NotificationChannel[userId]; ok {
		channel <- dto.NotificationData{
			ID:        notification.ID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		}
	}

	return nil
}
