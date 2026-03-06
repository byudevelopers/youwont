package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SupabaseID string             `bson:"supabase_id"   json:"supabase_id"`
	Name       string             `bson:"name"          json:"name"`
	Username   string             `bson:"username"      json:"username"`
	AvatarURL  *string            `bson:"avatar_url"    json:"avatar_url"`
	Points     int                `bson:"points"        json:"points"`
	CreatedAt  time.Time          `bson:"created_at"    json:"created_at"`
}
