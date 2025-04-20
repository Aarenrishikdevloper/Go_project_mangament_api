package task

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rishik.com/db"
	"rishik.com/enums"
	"rishik.com/models"
	"rishik.com/utils"
)

type TaskRequestBody struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	AssignedTo  *string `json:"assignedTo,omitempty"`
	DueDate     *string `json:"dueDate,omitempty"`
}

func CreateTaskService(workspaceId string, projectId string, userId string, body TaskRequestBody) (*models.Task, error) {
	projectObjId, err := primitive.ObjectIDFromHex(projectId)
	ctx := context.Background()
	if err != nil {
		return nil, errors.New("invalid project ID")
	}
	workspaceid, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}
	var project models.Project
	projectCollection := db.Database.Collection("projects")
	memberCollection := db.Database.Collection("members")
	err = projectCollection.FindOne(ctx, bson.M{
		"_id":       projectObjId,
		"workspace": workspaceid,
	}).Decode(&project)
	if err != nil {
		return nil, errors.New("project not found or does not belong to this workspace")
	}
	if body.AssignedTo != nil {
		assigneduserId, err := primitive.ObjectIDFromHex(*body.AssignedTo)
		if err != nil {
			return nil, errors.New("invalid assigned user ID")
		}
		count, err := memberCollection.CountDocuments(ctx, bson.M{
			"userId":      assigneduserId,
			"workspaceId": workspaceid,
		})
		if err != nil || count == 0 {
			return nil, errors.New("assigned user is not a member of this workspace")
		}

	}
	var dueDate *time.Time
	if body.DueDate != nil {
		parsedDuedate, err := time.Parse(time.RFC3339, *body.DueDate)
		if err == nil {
			dueDate = &parsedDuedate
		}
	}
	taskcode := utils.GenerateTaskCode()

	userIdobj, err := primitive.ObjectIDFromHex(*body.AssignedTo)
	if err != nil {
		return nil, errors.New("assigned user is not a member of this workspace")

	}
	taskAssigner, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       body.Title,
		Description: body.Description,
		Priority:    enums.TaskPriorityEnum(body.Priority),
		Status:      enums.TaskSTatusEnum(body.Status),
		AssignedTo:  &userIdobj,
		CreatedBy:   taskAssigner,
		WorkspaceID: workspaceid,
		ProjectID:   projectObjId,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		TaskCode:    taskcode,
	}

	_, err = db.Database.Collection("tasks").InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	return &task, nil

}

type UpTaskRequest struct {
	Title       string  `bson:"title"`
	Description *string `bson:"description,omitempty"`
	Priority    string  `bson:"priority"`
	Status      string  `bson:"status"`
	AssignedTo  *string `bson:"assignedTo,omitempty"`
	DueDate     *string `bson:"dueDate,omitempty"`
}

func UpdateTask(workspacId, projectid, taskId primitive.ObjectID, body UpTaskRequest) (*models.Task, error) {
	projectCollection := db.Database.Collection("projects")
	taskCollection := db.Database.Collection("tasks")
	ctx := context.Background()
	var project *models.Project
	err := projectCollection.FindOne(ctx, bson.M{
		"_id":       projectid,
		"workspace": workspacId,
	}).Decode(&project)
	if err != nil {
		return nil, fmt.Errorf("project not found or does not belong to this workspace: %w", err)

	}
	var task models.Task
	err = taskCollection.FindOne(ctx, bson.M{
		"_id":     taskId,
		"project": projectid,
	}).Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("task not found or does not belong to this project: %w", err)
	}
	update := bson.M{
		"$set": body,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updateTask models.Task
	err = taskCollection.FindOneAndUpdate(ctx, bson.M{
		"_id": taskId,
	}, update, opts).Decode(&updateTask)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)

	}
	return &updateTask, nil

}

var (
	ErrTastNotFound = errors.New("no task found")
)

func DeleteTaskService(workspaceId, taskId string) error {

	taskCollection := db.Database.Collection("tasks")
	ctx := context.Background()
	taskobjectId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}
	workspaceObjId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id":       taskobjectId,
		"workspace": workspaceObjId,
	}
	result, err := taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrTastNotFound
	}

	return nil

}

func GetTaskById(workspaceId, projectId, taskId string) (bson.M, error) {
	taskCollection := db.Database.Collection("tasks")
	objectTaskId, err := primitive.ObjectIDFromHex(taskId)
	ctx := context.Background()
	if err != nil {
		return nil, fmt.Errorf("invalid task ID")
	}
	objectProjectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return nil, fmt.Errorf("invalid Project ID")
	}
	objectWorkspaceId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("invalid Workspace ID")
	}
	projectCollection := db.Database.Collection("projects")
	projectFilter := bson.M{
		"_id":       objectProjectId,
		"workspace": objectWorkspaceId,
	}
	projectCount, err := projectCollection.CountDocuments(ctx, projectFilter)
	if err != nil {
		return nil, err
	}
	if projectCount == 0 {
		return nil, fmt.Errorf("project does not exist ")
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id":       objectTaskId,
			"project":   objectProjectId,
			"workspace": objectWorkspaceId,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "assignedTo",
			"foreignField": "_id",
			"as":           "assignedTo",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$assignedTo",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"title":                     1,
			"description":               1,
			"status":                    1,
			"priority":                  1,
			"workspace":                 1,
			"project":                   1,
			"assignedTo._id":            1,
			"assignedTo.name":           1,
			"assignedTo.profilePicture": 1,
			// explicitly exclude password
		}}},
	}
	cursor, err := taskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return results[0], nil

}

type TaskFilter struct {
	ProjectID  *string
	Status     []string
	Priority   []string
	AssignedTo []string
	Keyword    *string
	DueDate    *string
}

type Pagination struct {
	PageSize   int
	PageNumber int
}

type TaskResult struct {
	Tasks      []bson.M
	Pagination PaginationResult
}

type PaginationResult struct {
	PageSize   int
	PageNumber int
	TotalCount int64
	TotalPages int64
	Skip       int
}

func GetTaskServices(workSpaceId string, filters TaskFilter, pagination Pagination) (*TaskResult, error) {
	objectWorkspaceId, err := primitive.ObjectIDFromHex(workSpaceId)
	fmt.Print(workSpaceId)
	if err != nil {
		return nil, fmt.Errorf("invalid Workspace ID")
	}

	matchStage := bson.M{
		"workspace": objectWorkspaceId,
	}
	if filters.ProjectID != nil {
		objectProjectId, err := primitive.ObjectIDFromHex(*filters.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("invalid Workspace ID")
		}
		matchStage["project"] = objectProjectId
	}
	if len(filters.Status) > 0 {
		matchStage["status"] = bson.M{
			"$in": filters.Status,
		}
	}
	if len(filters.Priority) > 0 {
		matchStage["priority"] = bson.M{"$in": filters.Priority}
	}

	if len(filters.AssignedTo) > 0 {
		// Convert string IDs to ObjectIDs
		var assignedToIDs []primitive.ObjectID
		for _, id := range filters.AssignedTo {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				log.Print(err)
				return nil, err

			}
			assignedToIDs = append(assignedToIDs, objID)
		}
		matchStage["assignedTo"] = bson.M{"$in": assignedToIDs}
	}

	if filters.Keyword != nil && *filters.Keyword != "" {
		matchStage["title"] = bson.M{"$regex": primitive.Regex{Pattern: *filters.Keyword, Options: "i"}}
	}

	if filters.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *filters.DueDate)
		if err != nil {
			return nil, err
		}
		matchStage["dueDate"] = bson.M{"$eq": dueDate}
	}
	skip := (pagination.PageNumber - 1) * pagination.PageSize
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: pagination.PageSize}},
		// Lookup for assignedTo (user) information
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "users",
				"localField":   "assignedTo",
				"foreignField": "_id",
				"as":           "assignedToInfo",
			}},
		},
		// Unwind assignedToInfo (assuming it's a single user)
		{
			{Key: "$unwind", Value: bson.M{
				"path":                       "$assignedToInfo",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		// Project to include only needed fields from user
		{
			{Key: "$project", Value: bson.M{
				"assignedToInfo.password": 0,
			}},
		},
		// Lookup for project information
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "projects",
				"localField":   "project",
				"foreignField": "_id",
				"as":           "projectInfo",
			}},
		},
		// Unwind projectInfo
		{
			{Key: "$unwind", Value: bson.M{
				"path":                       "$projectInfo",
				"preserveNullAndEmptyArrays": true,
			}},
		},
	}
	ctx := context.Background()
	cursor, err := db.Database.Collection("tasks").Aggregate(ctx, pipeline)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var task []bson.M
	if err = cursor.All(ctx, &task); err != nil {
		return nil, err
	}
	totalCount, err := db.Database.Collection("tasks").CountDocuments(ctx, matchStage)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	totalPages := totalCount / int64(pagination.PageNumber)
	if totalCount%int64(pagination.PageSize) > 0 {
		totalPages++
	}
	return &TaskResult{
		Tasks: task,
		Pagination: PaginationResult{
			PageSize:   pagination.PageSize,
			PageNumber: pagination.PageNumber,
			TotalCount: totalCount,
			TotalPages: totalPages,
			Skip:       skip,
		},
	}, nil
}
