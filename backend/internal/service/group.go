package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"youwont.api/internal/model"
)

type GroupRepository interface {
	Create(ctx context.Context, group *model.Group) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Group, error)
	FindByMemberUserID(ctx context.Context, userID primitive.ObjectID) ([]model.Group, error)
	FindByInviteCode(ctx context.Context, code string) (*model.Group, error)
	PushMember(ctx context.Context, groupID primitive.ObjectID, member model.Member) error
}

type BetCounter interface {
	CountByGroupID(ctx context.Context, groupID primitive.ObjectID, status *string) (int64, error)
}

type GroupService struct {
	groups GroupRepository
	bets   BetCounter
}

func NewGroupService(groups GroupRepository, bets BetCounter) *GroupService {
	return &GroupService{groups: groups, bets: bets}
}

func (s *GroupService) Create(ctx context.Context, user *model.User, name, description string) (*model.Group, error) {
	code, err := generateInviteCode(6)
	if err != nil {
		return nil, err
	}

	group := &model.Group{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		InviteCode:  code,
		CreatedBy:   user.ID,
		Members: []model.Member{
			{UserID: user.ID, Role: "ADMIN", JoinedAt: time.Now()},
		},
		CreatedAt: time.Now(),
	}

	if err := s.groups.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) List(ctx context.Context, userID primitive.ObjectID) ([]model.Group, error) {
	return s.groups.FindByMemberUserID(ctx, userID)
}

func (s *GroupService) Get(ctx context.Context, groupID, userID primitive.ObjectID) (*model.Group, error) {
	group, err := s.groups.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrNotFound
	}
	if !isMember(group, userID) {
		return nil, ErrForbidden
	}
	return group, nil
}

func (s *GroupService) JoinByCode(ctx context.Context, user *model.User, code string) (*model.Group, error) {
	group, err := s.groups.FindByInviteCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrNotFound
	}
	if isMember(group, user.ID) {
		return nil, ErrAlreadyExists
	}

	member := model.Member{
		UserID:   user.ID,
		Role:     "MEMBER",
		JoinedAt: time.Now(),
	}
	if err := s.groups.PushMember(ctx, group.ID, member); err != nil {
		return nil, err
	}

	group.Members = append(group.Members, member)
	return group, nil
}

func isMember(group *model.Group, userID primitive.ObjectID) bool {
	for _, m := range group.Members {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

func generateInviteCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}
