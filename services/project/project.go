package project

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	//"rishik.com/controller/taskcontroller"
	"rishik.com/db"
	"rishik.com/enums"
	"rishik.com/models"
)

type ProjectInput struct {
	Emoji       *string `json:"emoji"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
type ProjectSummary struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Emoji       string             `bson:"emoji,omitempty" json:"emoji,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
}

func CreateProject(userId, workspaceId string, body ProjectInput) (*models.Project, error) {
	objId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	userobjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		Name:        body.Name,
		WorkspaceID: objId,
		CreatedBy:   userobjId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if body.Emoji != nil {
		project.Emoji = *body.Emoji
	}
	if body.Description != nil {
		project.Description = body.Description
	}
	collection := db.Database.Collection("projects")
	res, err := collection.InsertOne(ctx, project)

	if err != nil {
		return nil, err
	}
	project.ID = res.InsertedID.(primitive.ObjectID)
	return project, nil
}

// Check if the user is a member of the workspace
type UpdateProjectBody struct {
	Emoji       *string `json:"emoji,omitempty"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func UpdatePeoject(workspaceId, projectId string, body UpdateProjectBody) (*models.Project, error) {
	ctx := context.Background()
	projectCollection := db.Database.Collection("projects")
	wID, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, errors.New("invalid workspace ID")
	}
	pId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, errors.New("invalid workspace ID")
	}
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{
		"_id":       pId,
		"workspace": wID,
	}).Decode(&project)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("project not found")
	} else if err != nil {
		return nil, err
	}
	update := bson.M{
		"name": body.Name,
	}
	if body.Emoji != nil {
		update["emoji"] = *body.Emoji
	}
	if body.Description != nil {
		update["description"] = *body.Description
	}
	_, err = projectCollection.UpdateOne(ctx, bson.M{
		"_id": pId,
	},
		bson.M{
			"$set": update,
		})
	if err != nil {
		return nil, err
	}
	err = projectCollection.FindOne(ctx, bson.M{"_id": pId}).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
func DeleteprojectService(workspaceId, projectId string) (*models.Project, error) {
	ctx := context.Background()
	projectCollection := db.Database.Collection("projects")
	taskcollection := db.Database.Collection("tasks")
	wId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, errors.New("invalid workspace ID")
	}
	pId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}
	var project models.Project
	err = projectCollection.FindOne(ctx, bson.M{
		"_id":       pId,
		"workspace": wId,
	}).Decode(&project)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("project not found or does not belong to the specified workspace")
	} else if err != nil {
		return nil, err
	}
	_, err = projectCollection.DeleteOne(ctx, bson.M{"_id": pId})
	if err != nil {
		return nil, err
	}
	_, err = taskcollection.DeleteMany(ctx, bson.M{"project": pId})
	if err != nil {
		return nil, err
	}
	return &project, nil

}

func GetProjectBYidAndWorkspaceIdService(workspaceId string, projectId string) (*ProjectSummary, error) {
	projectCollection := db.Database.Collection("projects")
	ctx := context.Background()
	wId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, errors.New("invalid workspaceid")
	}
	pid, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	var project ProjectSummary
	err = projectCollection.FindOne(ctx, bson.M{
		"_id":       pid,
		"workspace": wId,
	}).Decode(&project)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("project not found or does not belong to the specified workspace")
	} else if err != nil {
		return nil, err
	}
	return &project, nil

}

type User struct {
	ID             string `bson:"_id,omitempty" json:"_id"`
	Name           string `bson:"name" json:"name"`
	ProfilePicture string `bson:"profilePicture" json:"profilePicture"`
}

type Project struct {
	ID          string    `bson:"_id,omitempty" json:"_id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Emoji       string    `bson:"emoji" json:"emoji"`
	Workspace   string    `bson:"workspace" json:"workspace"`
	CreatedBy   []User    `bson:"createdBy" json:"createdBy"` // It's an array due to aggregation `$lookup`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}

type PaginationProjects struct {
	Projects   []Project `json:"projects"`
	TotalCount int64     `json:"totalCount"`
	TotalPages int       `json:"totalPages"`
	Skip       int       `json:"skip"`
}

func GetProjectsInWorkspaces(workspaceId string, pageSize int64, pageNumber int64) (*PaginationProjects, error) {
	projectCollection := db.Database.Collection("projects")
	ctx := context.Background()
	wId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, errors.New("invalid workspace ID")
	}
	totalCount, err := projectCollection.CountDocuments(ctx, bson.M{"workspace": wId})
	if err != nil {
		return nil, err
	}
	skip := (pageNumber - 2) * pageSize
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"workspace": wId}}},
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: pageSize}},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "users",
				"localField":   "createdBy",
				"foreignField": "_id",
				"as":           "createdBy",
			}},
		},
		{
			{Key: "$project", Value: bson.M{
				"name":        1,
				"workspace":   1,
				"emoji":       1,
				"description": 1,
				"createdAt":   1,
				"createdBy": bson.M{
					"_id":            1,
					"name":           1,
					"profilePicture": 1,
				},
			}},
		},
	}
	cursor, err := projectCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var projects []Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return &PaginationProjects{
		Projects:   projects,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Skip:       int(skip),
	}, nil
}

type ProjectAnalytics struct {
	TotalTask     int64 `json:"totalTasks"`
	OverdueTask   int64 `json:"overdueTasks"`
	CompletedTask int64 `json:"completedTasks"`
}

func GetProjectAnalytics(workspaceId, projectId primitive.ObjectID) (*ProjectAnalytics, error) {
	projectCollection := db.Database.Collection("projects")
	taskCollection := db.Database.Collection("tasks")
	ctx := context.Background()
	var project models.Project
	var err error

	err = projectCollection.FindOne(ctx, bson.M{"_id": projectId}).Decode(&project)
	if err != nil || project.WorkspaceID != workspaceId {
		return nil, errors.New("project not found or does not belong to this workspace")

	}
	currentDate := time.Now()
	totalTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
		"project":   projectId,
	})
	if err != nil {
		return nil, err
	}
	overdueTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
		"project":   projectId,
		"dueDate":   bson.M{"$lt": currentDate},
		"status":    bson.M{"$ne": enums.Done},
	})
	if err != nil {
		return nil, err
	}
	completedTask, err := taskCollection.CountDocuments(ctx, bson.M{
		"workspace": workspaceId,
		"status":    enums.Done,
		"project":   projectId,
	})
	if err != nil {
		return nil, err
	}
	analysis := &ProjectAnalytics{
		TotalTask:     totalTask,
		CompletedTask: completedTask,
		OverdueTask:   overdueTask,
	}
	return analysis, nil

}
