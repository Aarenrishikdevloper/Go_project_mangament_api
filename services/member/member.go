package member

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"rishik.com/db"
	"rishik.com/enums"
	"rishik.com/models"
)

type JoinWorkspaceResponse struct {
	WorkspaceID string
	Role        string
}

func JoinWorkspace(userId primitive.ObjectID, invitecode string) (*JoinWorkspaceResponse, error) {
	workspaceColl := db.Database.Collection("workspaces")
	memberColl := db.Database.Collection("members")
	roleColl := db.Database.Collection("roles")
	ctx := context.Background()
	var workspaces models.Workspace
	err := workspaceColl.FindOne(ctx, bson.M{"inviteCode": invitecode}).Decode(&workspaces)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid invite code or workspace not found")
		}
		return nil, err
	}
	filter := bson.M{"userId": userId, "workspaceId": workspaces.ID}
	var existingMember models.Member
	err = memberColl.FindOne(ctx, filter).Decode(&existingMember)
	if err == nil {
		return nil, errors.New("you are already a member")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}
	var role models.Role
	err = roleColl.FindOne(ctx, bson.M{"name": enums.Member}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	newMember := models.Member{
		UserID:      userId,
		WorkspaceID: workspaces.ID,
		RoleID:      role.ID,
		JoinedAt:    time.Now(),
	}
	_, err = memberColl.InsertOne(ctx, newMember)
	if err != nil {
		return nil, err
	}
	return &JoinWorkspaceResponse{
		WorkspaceID: workspaces.ID.Hex(),
		Role:        string(role.Name),
	}, nil
}
