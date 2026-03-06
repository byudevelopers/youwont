package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"youwont.api/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	CreateMany(ctx context.Context, notifs []model.Notification) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]model.Notification, error)
	CountUnread(ctx context.Context, userID primitive.ObjectID) (int64, error)
	MarkRead(ctx context.Context, id, userID primitive.ObjectID) error
	MarkAllRead(ctx context.Context, userID primitive.ObjectID) (int64, error)
}

type NotificationService struct {
	notifs NotificationRepository
}

func NewNotificationService(notifs NotificationRepository) *NotificationService {
	return &NotificationService{notifs: notifs}
}

func (s *NotificationService) List(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]model.Notification, bool, error) {
	notifs, err := s.notifs.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if int64(len(notifs)) > limit {
		hasMore = true
		notifs = notifs[:limit]
	}

	return notifs, hasMore, nil
}

func (s *NotificationService) UnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return s.notifs.CountUnread(ctx, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, notifID, userID primitive.ObjectID) error {
	return s.notifs.MarkRead(ctx, notifID, userID)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return s.notifs.MarkAllRead(ctx, userID)
}
