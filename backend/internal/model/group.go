package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	UserID   primitive.ObjectID `bson:"user_id"   json:"user_id"`
	Role     string             `bson:"role"      json:"role"`
	JoinedAt time.Time          `bson:"joined_at" json:"joined_at"`
}

type Group struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name"          json:"name"`
	Description string             `bson:"description"   json:"description"`
	InviteCode  string             `bson:"invite_code"   json:"invite_code"`
	CreatedBy   primitive.ObjectID `bson:"created_by"    json:"created_by"`
	Members     []Member           `bson:"members"       json:"members"`
	CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
}
