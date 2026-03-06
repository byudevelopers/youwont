package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"youwont.api/internal/model"
)

type GroupRepo struct {
	col *mongo.Collection
}

func NewGroupRepo(db *mongo.Database) *GroupRepo {
	col := db.Collection("groups")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "invite_code", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "members.user_id", Value: 1}}},
	})
	return &GroupRepo{col: col}
}

func (r *GroupRepo) Create(ctx context.Context, group *model.Group) error {
	_, err := r.col.InsertOne(ctx, group)
	return err
}

func (r *GroupRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Group, error) {
	var group model.Group
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &group, err
}

func (r *GroupRepo) FindByMemberUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Group, error) {
	cursor, err := r.col.Find(ctx, bson.M{"members.user_id": userID})
	if err != nil {
		return nil, err
	}
	var groups []model.Group
	return groups, cursor.All(ctx, &groups)
}

func (r *GroupRepo) FindByInviteCode(ctx context.Context, code string) (*model.Group, error) {
	var group model.Group
	err := r.col.FindOne(ctx, bson.M{"invite_code": code}).Decode(&group)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &group, err
}

func (r *GroupRepo) PushMember(ctx context.Context, groupID primitive.ObjectID, member model.Member) error {
	_, err := r.col.UpdateByID(ctx, groupID, bson.M{
		"$push": bson.M{"members": member},
	})
	return err
}
