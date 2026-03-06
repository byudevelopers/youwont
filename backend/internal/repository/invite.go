package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"youwont.api/internal/model"
)

type InviteRepo struct {
	col *mongo.Collection
}

func NewInviteRepo(db *mongo.Database) *InviteRepo {
	col := db.Collection("invites")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "invitee_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "group_id", Value: 1}, {Key: "invitee_id", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return &InviteRepo{col: col}
}

func (r *InviteRepo) Create(ctx context.Context, invite *model.Invite) error {
	_, err := r.col.InsertOne(ctx, invite)
	return err
}

func (r *InviteRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Invite, error) {
	var invite model.Invite
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&invite)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &invite, err
}

func (r *InviteRepo) FindPendingByInvitee(ctx context.Context, inviteeID primitive.ObjectID) ([]model.Invite, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.col.Find(ctx, bson.M{"invitee_id": inviteeID, "status": "PENDING"}, opts)
	if err != nil {
		return nil, err
	}
	var invites []model.Invite
	return invites, cursor.All(ctx, &invites)
}

func (r *InviteRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"status": status}})
	return err
}
