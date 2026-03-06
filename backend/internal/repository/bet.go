package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"youwont.api/internal/model"
)

type BetRepo struct {
	col *mongo.Collection
}

func NewBetRepo(db *mongo.Database) *BetRepo {
	col := db.Collection("bets")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "group_id", Value: 1}}},
		{Keys: bson.D{{Key: "wagers.user_id", Value: 1}}},
	})
	return &BetRepo{col: col}
}

func (r *BetRepo) Create(ctx context.Context, bet *model.Bet) error {
	_, err := r.col.InsertOne(ctx, bet)
	return err
}

func (r *BetRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Bet, error) {
	var bet model.Bet
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&bet)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &bet, err
}

func (r *BetRepo) FindByGroupID(ctx context.Context, groupID primitive.ObjectID, status *string) ([]model.Bet, error) {
	filter := bson.M{"group_id": groupID}
	if status != nil {
		filter["status"] = *status
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var bets []model.Bet
	return bets, cursor.All(ctx, &bets)
}

func (r *BetRepo) CountByGroupID(ctx context.Context, groupID primitive.ObjectID, status *string) (int64, error) {
	filter := bson.M{"group_id": groupID}
	if status != nil {
		filter["status"] = *status
	}
	return r.col.CountDocuments(ctx, filter)
}

func (r *BetRepo) HasOpenBetsByUser(ctx context.Context, groupID, userID primitive.ObjectID) (bool, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{
		"group_id": groupID,
		"status":   "OPEN",
		"$or": bson.A{
			bson.M{"creator_id": userID},
			bson.M{"decider_id": userID},
			bson.M{"wagers.user_id": userID},
		},
	})
	return count > 0, err
}

func (r *BetRepo) PushWager(ctx context.Context, betID primitive.ObjectID, wager model.Wager) error {
	_, err := r.col.UpdateByID(ctx, betID, bson.M{
		"$push": bson.M{"wagers": wager},
	})
	return err
}

func (r *BetRepo) Resolve(ctx context.Context, betID primitive.ObjectID, winningSide string, resolvedAt time.Time) error {
	_, err := r.col.UpdateByID(ctx, betID, bson.M{
		"$set": bson.M{
			"status":       "RESOLVED",
			"winning_side": winningSide,
			"resolved_at":  resolvedAt,
		},
	})
	return err
}

func (r *BetRepo) UpdateDecider(ctx context.Context, betID, deciderID primitive.ObjectID) error {
	_, err := r.col.UpdateByID(ctx, betID, bson.M{
		"$set": bson.M{"decider_id": deciderID},
	})
	return err
}
