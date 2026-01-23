package usecase

import (
	"context"
	"errors"

	"github.com/Ramsi97/edu-social-backend/internal/group/domain"
	"github.com/Ramsi97/edu-social-backend/internal/group/repository/interfaces"
	"github.com/google/uuid"
)

type groupChatUseCase struct {
	repo      interfaces.GroupChatRepo
}

func NewGroupChatUseCase(repo interfaces.GroupChatRepo) domain.GroupChatUseCase {
	return &groupChatUseCase{
		repo:      repo,
	}
}

func (g *groupChatUseCase) CreateGroup(ctx context.Context, ownerID uuid.UUID, groupName string) (uuid.UUID, error) {
	_, err := g.repo.GetGroup(ctx, groupName) 
	if err == nil {
        return uuid.Nil, errors.New("group already exists")
    }

	groupID := uuid.New()

	group := domain.Group{
		ID:      groupID,
        Name:    groupName,
        OwnerID: ownerID,
	}

	if err := g.repo.CreateGroup(ctx, &group); err != nil {
        return uuid.Nil, err
    }

    return groupID, nil
}

func (g *groupChatUseCase) GetMessages(ctx context.Context, groupID uuid.UUID, limit int) ([]*domain.Message, error) {
	return g.repo.GetMessages(ctx, groupID, limit)
}

func (g *groupChatUseCase) JoinGroup(ctx context.Context, groupName string, userID uuid.UUID) error {
	groupID, err := g.repo.GetGroup(ctx, groupName)
    if err != nil {
        return errors.New("group didn't exist")
    }

    return g.repo.JoinGroup(ctx, groupID, userID)
}

func (g *groupChatUseCase) LeaveGroup(ctx context.Context, groupName string, userID uuid.UUID) error {
	groupID, err := g.repo.GetGroup(ctx, groupName)

	if err != nil {
		return errors.New("group didn't exist")
	}
	return g.repo.LeaveGroup(ctx, groupID, userID)
}

func (g *groupChatUseCase) SendMessage(ctx context.Context, msg *domain.Message) error {
	member, err := g.repo.IsMember(ctx, msg.AuthorID, msg.GroupID)

	if err != nil {
		return err
	}

	if !member{
		return errors.New("not allowed to post to this group")
	}

	if err := g.repo.SaveMessage(ctx, msg); err != nil {
        return err
    }

	return nil
}


func (uc *groupChatUseCase) GetGroupsForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Group, error) {
	// 1️⃣ Validate input (optional)
	if userID == uuid.Nil {
		return nil, domain.ErrGroupNotFound // or another domain error for invalid user
	}

	// 2️⃣ Call repository
	groups, err := uc.repo.GetGroupsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}