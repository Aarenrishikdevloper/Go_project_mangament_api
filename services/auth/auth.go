package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rishik.com/db"
	"rishik.com/enums"
	"rishik.com/models"
	"rishik.com/utils"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func ptr(s string) *string {
	return &s
}

func RegisterUserService(body RegisterRequest) (*models.User, error) {
	ctx := context.Background()
	session, err := db.Client.StartSession()
	if err != nil {
		return nil, err
	}
	var registeredUser *models.User
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		UserCollection := db.Database.Collection("users")
		accountCollection := db.Database.Collection("accounts")
		workspaceCollection := db.Database.Collection("workspaces")
		rolesCollection := db.Database.Collection("roles")
		memberCollection := db.Database.Collection("members")
		var existingUser models.User
		err := UserCollection.FindOne(sc, bson.M{"email": body.Email}).Decode(&existingUser)
		if err == nil {
			return nil, errors.New("email already exists")
		} else if err != mongo.ErrNoDocuments {
			return nil, err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedpasswordStr := string(hashedPassword)
		user := models.User{
			Email:     body.Email,
			Name:      &body.Name,
			Password:  &hashedpasswordStr,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRes, err := UserCollection.InsertOne(sc, user)
		if err != nil {
			return nil, err
		}
		userID := userRes.InsertedID.(primitive.ObjectID)
		registeredUser = &models.User{
			ID:        userID,
			Email:     body.Email,
			Name:      &body.Name,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		account := models.Account{
			UserID:     userID,
			Provider:   enums.Email,
			ProviderID: body.Email,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = accountCollection.InsertOne(sc, account)
		if err != nil {
			return nil, err
		}
		workspace := models.Workspace{
			Name:        "My Workspace",
			Description: ptr("Workspace created for" + body.Name),
			OwnerID:     userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			InviteCode:  utils.GenerateInviteCode(),
		}
		workspaceRes, err := workspaceCollection.InsertOne(sc, workspace)
		if err != nil {
			return nil, err
		}
		worspaceId := workspaceRes.InsertedID.(primitive.ObjectID)
		var ownerRole models.Role
		err = rolesCollection.FindOne(sc, bson.M{"name": enums.Owner}).Decode(&ownerRole)
		if err != nil {
			return nil, errors.New("owner role not found")
		}
		member := models.Member{
			UserID:      userID,
			WorkspaceID: worspaceId,
			RoleID:      ownerRole.ID,
			JoinedAt:    time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err = memberCollection.InsertOne(sc, member)
		if err != nil {
			return nil, err
		}
		_, err = UserCollection.UpdateOne(sc, bson.M{"_id": userID}, bson.M{"$set": bson.M{
			"currentWorkspace": worspaceId,
			"updatedAt":        time.Now(),
		}})
		if err != nil {
			return err, nil
		}
		registeredUser.CurrentWorkspace = &worspaceId
		return nil, nil
	})
	return registeredUser, nil
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	//Provider enums.ProviderEnum
	Password string `json:"password" validate:"required,min=8"`
}

func LoginService(body LoginRequest) (*models.User, error) {
	ctx := context.Background()
	accountCollection := db.Database.Collection("accounts")
	var account models.Account
	err := accountCollection.FindOne(ctx, bson.M{
		"provider":   enums.Email,
		"providerId": body.Email,
	}).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, errors.New("invalid email or password")

		}

	}
	var user models.User
	UserCollection := db.Database.Collection("users")
	err = UserCollection.FindOne(ctx, bson.M{
		"_id": account.UserID,
	}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, errors.New("invalid email or password")

		}

	}
	err = user.ComparePassword(body.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user.OmitPassword(), nil

}
