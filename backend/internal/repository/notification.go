package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"youwont.api/internal/model"
)

type NotificationRepo struct {
	col *mongo.Collection
}

func NewNotificationRepo(db *mongo.Database) *NotificationRepo {
	col := db.Collection("notifications")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "read", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}},
	})
	return &NotificationRepo{col: col}
}

func (r *NotificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	_, err := r.col.InsertOne(ctx, notif)
	return err
}

func (r *NotificationRepo) CreateMany(ctx context.Context, notifs []model.Notification) error {
	docs := make([]interface{}, len(notifs))
	for i, n := range notifs {
		docs[i] = n
	}
	_, err := r.col.InsertMany(ctx, docs)
	return err
}

func (r *NotificationRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]model.Notification, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(page * limit).
		SetLimit(limit + 1)
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	var notifs []model.Notification
	return notifs, cursor.All(ctx, &notifs)
}

func (r *NotificationRepo) CountUnread(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{"user_id": userID, "read": false})
}

func (r *NotificationRepo) MarkRead(ctx context.Context, id, userID primitive.ObjectID) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": bson.M{"read": true}})
	return err
}

func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	result, err := r.col.UpdateMany(ctx, bson.M{"user_id": userID, "read": false}, bson.M{"$set": bson.M{"read": true}})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
