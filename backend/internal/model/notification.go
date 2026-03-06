package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id"       json:"user_id"`
	Type      string             `bson:"type"          json:"type"`
	RefType   string             `bson:"ref_type"      json:"ref_type"`
	RefID     primitive.ObjectID `bson:"ref_id"        json:"ref_id"`
	Message   string             `bson:"message"       json:"message"`
	Read      bool               `bson:"read"          json:"read"`
	CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
}
