package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wager struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   primitive.ObjectID `bson:"user_id"       json:"user_id"`
	Side     string             `bson:"side"          json:"side"`
	Amount   int                `bson:"amount"        json:"amount"`
	PlacedAt time.Time          `bson:"placed_at"     json:"placed_at"`
}

type Bet struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID     primitive.ObjectID `bson:"group_id"      json:"group_id"`
	Title       string             `bson:"title"         json:"title"`
	Description string             `bson:"description"   json:"description"`
	CreatorID   primitive.ObjectID `bson:"creator_id"    json:"creator_id"`
	DeciderID   primitive.ObjectID `bson:"decider_id"    json:"decider_id"`
	EndDate     time.Time          `bson:"end_date"      json:"end_date"`
	Status      string             `bson:"status"        json:"status"`
	WinningSide *string            `bson:"winning_side"  json:"winning_side"`
	Wagers      []Wager            `bson:"wagers"        json:"wagers"`
	ResolvedAt  *time.Time         `bson:"resolved_at"   json:"resolved_at"`
	CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
}
