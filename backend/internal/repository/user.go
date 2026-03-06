package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"youwont.api/internal/model"
)

type UserRepo struct {
	col *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	col := db.Collection("users")
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "supabase_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return &UserRepo{col: col}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) FindBySupabaseID(ctx context.Context, supabaseID string) (*model.User, error) {
	var user model.User
	err := r.col.FindOne(ctx, bson.M{"supabase_id": supabaseID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepo) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]model.User, error) {
	cursor, err := r.col.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var users []model.User
	return users, cursor.All(ctx, &users)
}

func (r *UserRepo) SearchByUsername(ctx context.Context, query string, limit int64) ([]model.User, error) {
	filter := bson.M{"username": bson.M{"$regex": "^" + query, "$options": "i"}}
	opts := options.Find().SetLimit(limit)
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var users []model.User
	return users, cursor.All(ctx, &users)
}

func (r *UserRepo) AddPoints(ctx context.Context, userID primitive.ObjectID, amount int) error {
	_, err := r.col.UpdateByID(ctx, userID, bson.M{"$inc": bson.M{"points": amount}})
	return err
}

func (r *UserRepo) DeductPoints(ctx context.Context, userID primitive.ObjectID, amount int) error {
	return r.AddPoints(ctx, userID, -amount)
}
