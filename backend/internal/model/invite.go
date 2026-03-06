package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invite struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"   json:"id"`
	GroupID       primitive.ObjectID `bson:"group_id"        json:"group_id"`
	GroupName     string             `bson:"group_name"      json:"group_name"`
	InvitedBy     primitive.ObjectID `bson:"invited_by"      json:"invited_by"`
	InvitedByName string             `bson:"invited_by_name" json:"invited_by_name"`
	InviteeID     primitive.ObjectID `bson:"invitee_id"      json:"invitee_id"`
	Status        string             `bson:"status"          json:"status"`
	CreatedAt     time.Time          `bson:"created_at"      json:"created_at"`
}
