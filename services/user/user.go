package user

import (
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rishik.com/db"
	"rishik.com/models"
)

func GetCurrentUsersService(userId string) (*models.User, error) {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid userid")
	}
	var user models.User
	UserCollection := db.Database.Collection("users")
	err = UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	//user.Password = ""
	return user.OmitPassword(), nil
}
