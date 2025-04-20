package workspace

import (
	"context"
	"errors"
	"fmt"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rishik.com/db"
	"rishik.com/enums"
	"rishik.com/models"
	"rishik.com/utils"
)

type WorkspaceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func CreateWorkSpaceService(body WorkspaceRequest, userId string) (*models.Workspace, error) {
	ctx := context.Background()
	UserCollection := db.Database.Collection("users")
	workspaceCollection := db.Database.Collection("workspaces")
	rolesCollection := db.Database.Collection("roles")
	memberCollection := db.Database.Collection("members")
	var user models.User

	err := UserCollection.FindOne(ctx, bson.M{
		"_id": userId,
	}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return nil, errors.New("user does not exist")

		}
	}

	user_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return nil, err
	}
	var ownerRole models.Role
	err = rolesCollection.FindOne(ctx, bson.M{"name": enums.Owner}).Decode(&ownerRole)
	if err != nil {
		return nil, errors.New("owner role not found")
	}
	workspace := models.Workspace{
		Name:        body.Name,
		Description: body.Description,
		OwnerID:     user_id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		InviteCode:  utils.GenerateInviteCode(),
	}
	fmt.Print(user.ID)
	workspaceRes, err := workspaceCollection.InsertOne(ctx, workspace)
	if err != nil {
		return nil, err
	}
	workspaceId := workspaceRes.InsertedID.(primitive.ObjectID)
	member := models.Member{
		UserID:      user_id,
		WorkspaceID: workspaceId,
		RoleID:      ownerRole.ID,
		JoinedAt:    time.Now(),
	}
	_, err = memberCollection.InsertOne(ctx, member)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	_, err = UserCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"currentWorkspace": workspaceId,
		"updatedAt":        time.Now(),
	}})
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &workspace, nil

}
func UpdateWorkspace(workspaceId string, name string, description *string) (*models.Workspace, error) {
	ctx := context.Background()
	workspaceCollection := db.Database.Collection("workspaces")
	objId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, err

	}
	var workspace models.Workspace
	err = workspaceCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&workspace)
	if errors.Is(err, mongo.ErrNilDocument) {
		return nil, errors.New("workspace does not exists")

	}
	update := bson.M{}
	if name != "" {
		update["name"] = name
	}
	if description != nil {
		update["description"] = *description

	}
	if len(update) > 0 {
		_, err := workspaceCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			return nil, err
		}
		if name != "" {
			workspace.Name = name
		}
		if description != nil {
			workspace.Description = description
		}

	}
	return &workspace, nil

}
func DeleteWorkspaceService(workspaceId, userid string) (*primitive.ObjectID, error) {
	ctx := context.Background()
	session, err := db.Client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	result, err := session.WithTransaction(ctx, func(sessctx mongo.SessionContext) (interface{}, error) {
		UserCollection := db.Database.Collection("users")
		workspaceCollection := db.Database.Collection("workspaces")
		projectcol := db.Database.Collection("projects")
		objId, err := primitive.ObjectIDFromHex(workspaceId)
		if err != nil {
			return nil, err

		}
		userobjId, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			return nil, err

		}
		taskcol := db.Database.Collection("tasks")
		memberCollection := db.Database.Collection("members")
		var workspace models.Workspace
		err = workspaceCollection.FindOne(sessctx, bson.M{
			"_id": objId,
		}).Decode(&workspace)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("workspace not found")
		} else if err != nil {
			return nil, err
		}
		var user models.User
		err = UserCollection.FindOne(sessctx, bson.M{"_id": userobjId}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		} else if err != nil {
			return nil, err
		}

		if workspace.OwnerID.Hex() != userid {
			return nil, fmt.Errorf("you are not authorized to delete this workspace")
		}
		_, err = projectcol.DeleteMany(sessctx, bson.M{"workspace": workspaceId})
		if err != nil {
			return nil, err
		}
		_, err = taskcol.DeleteMany(sessctx, bson.M{"workspace": workspaceId})
		if err != nil {
			return nil, err
		}
		_, err = memberCollection.DeleteMany(sessctx, bson.M{"workspace": workspaceId})
		if err != nil {
			return nil, err
		}
		if user.CurrentWorkspace != nil && user.CurrentWorkspace.Hex() == workspaceId {
			var member models.Member
			err := memberCollection.FindOne(sessctx, bson.M{"userId": userobjId}).Decode(&member)
			if err == mongo.ErrNoDocuments {
				user.CurrentWorkspace = nil
			} else if err != nil {
				return nil, err
			} else {
				user.CurrentWorkspace = &member.WorkspaceID
			}
			_, err = UserCollection.UpdateOne(
				sessctx,
				bson.M{"_id": userid},
				bson.M{"$set": bson.M{"currentWorkspace": user.CurrentWorkspace}},
			)
			if err != nil {
				return nil, err
			}

		}

		_, err = workspaceCollection.DeleteOne(sessctx, bson.M{"_id": objId})
		if err != nil {
			return nil, err
		}
		return user.CurrentWorkspace, nil

	})
	if err != nil {
		return nil, err
	}
	currentWorkspace := result.(*primitive.ObjectID)
	return currentWorkspace, nil
}
func GetAllWorkSpacesUserIsMemberService(userId string) ([]models.Workspace, error) {
	ctx := context.Background()
	userobjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err

	}
	memberCollection := db.Database.Collection("members")
	workspaceCollection := db.Database.Collection("workspaces")
	cursor, err := memberCollection.Find(ctx, bson.M{"userId": userobjId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var member []models.Member
	if err := cursor.All(ctx, &member); err != nil {
		return nil, err
	}
	var workSpaceIds []primitive.ObjectID
	for _, member := range member {
		workSpaceIds = append(workSpaceIds, member.WorkspaceID)
	}
	if len(workSpaceIds) == 0 {
		return []models.Workspace{}, nil
	}
	workspaceCursor, err := workspaceCollection.Find(ctx, bson.M{"_id": bson.M{"$in": workSpaceIds}})
	if err != nil {
		return nil, err
	}
	defer workspaceCursor.Close(ctx)
	var workspaces []models.Workspace
	if err := workspaceCursor.All(ctx, &workspaces); err != nil {
		return nil, err
	}
	return workspaces, nil

}

type workspaceWithMember struct {
	Workspace models.Workspace `json:"workspace"`
	Members   []models.Member  `json:"members"`
}

func GetWorkSpacebyIdService(workspaceId string) (*workspaceWithMember, error) {
	ctx := context.Background()
	fmt.Print(workspaceId)
	objId, err := primitive.ObjectIDFromHex(workspaceId)
	workspaceCollection := db.Database.Collection("workspaces")
	memberCollection := db.Database.Collection("members")
	rolesCollection := db.Database.Collection("roles")
	if err != nil {
		return nil, err
	}
	var workspace models.Workspace
	err = workspaceCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&workspace)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("workspace not found")
		}
		return nil, err
	}
	cursor, err := memberCollection.Find(ctx, bson.M{"workspaceId": objId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var members []models.Member
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	for i, member := range members {
		var role models.Role
		err = rolesCollection.FindOne(ctx, bson.M{"_id": member.RoleID}).Decode(&role)
		if err == nil {
			members[i].RoleID = role.ID
		}
	}
	result := &workspaceWithMember{
		Workspace: workspace,
		Members:   members,
	}
	return result, nil

}

type Analytics struct {
	TotalTask     int64 `json:"totalTasks"`
	OverdueTask   int64 `json:"overdueTasks"`
	CompletedTask int64 `json:"completedTasks"`
}

func GetWorkspaceAnalyticsServices(workspaceId primitive.ObjectID) (*Analytics, error) {
	taskCollection := db.Database.Collection("tasks")
	currentDate := time.Now()
	ctx := context.Background()
	totalTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
	})
	if err != nil {
		return nil, err
	}
	overdueTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
		"dueDate":   bson.M{"$lt": currentDate},
		"status":    bson.M{"$ne": enums.Done},
	})
	if err != nil {
		return nil, err
	}
	completedTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
		"status":    enums.Done,
	})
	if err != nil {
		return nil, err
	}
	analytics := &Analytics{
		TotalTask:     totalTask,
		OverdueTask:   overdueTask,
		CompletedTask: completedTask,
	}
	return analytics, nil
}
