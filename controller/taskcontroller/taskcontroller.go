package taskcontroller

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rishik.com/services/task"
)

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func CreateTask(c fiber.Ctx) error {
	var body task.TaskRequestBody
	workspaceId := c.Params("workspaceId")
	projectId := c.Params("projectId")
	userId := c.Locals("userId")
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	task, err := task.CreateTaskService(workspaceId, projectId, userId.(string), body)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to Create Task",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task Created Sucessfully",
		"task":    task,
	})
}
func UpdateTaskController(c fiber.Ctx) error {
	var body task.UpTaskRequest
	workspaceId := c.Params("workspaceId")
	projectId := c.Params("projectId")
	taskId := c.Params("taskId")

	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	projectid, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Project ID")
	}
	workspaceid, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Workspace ID")
	}
	taskid, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Workspace ID")
	}

	task, err := task.UpdateTask(workspaceid, projectid, taskid, body)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to Update Task",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task Update Sucessfully",
		"task":    task,
	})
}
func DeleteTaskController(c fiber.Ctx) error {
	worspaceId := c.Params("workspaceId")
	taskid := c.Params("taskId")
	err := task.DeleteTaskService(worspaceId, taskid)
	if err != nil {
		if err == task.ErrTastNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "No Task found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")

		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task delete Sucessfully",
	})
}
func GetTaskBYIdController(c fiber.Ctx) error {
	workspaceId := c.Params("workspaceId")
	projectId := c.Params("projectId")
	taskId := c.Params("taskId")
	task, err := task.GetTaskById(workspaceId, projectId, taskId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong")
	}
	return c.JSON(fiber.Map{
		"task": task,
	})
}
func GetAllWorkspace(c fiber.Ctx) error {
	workspaceId := c.Params("workspaceId")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace ID is required",
		})
	}
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("pageSize", "10")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	filters := task.TaskFilter{
		ProjectID:  ptr(c.Query("projectId")),
		Status:     strings.Split(c.Query("status", ""), ","), // Split comma-separated values into a slice
		Priority:   strings.Split(c.Query("priority", ""), ","),
		AssignedTo: strings.Split(c.Query("assignedTo", ""), ","),
		Keyword:    ptr(c.Query("keyword")),
		DueDate:    ptr(c.Query("dueDate")),
	}

	// Helper function to convert string to *string

	// Prepare pagination
	pagination := task.Pagination{
		PageSize:   int(pageSize),
		PageNumber: int(page),
	}
	result, err := task.GetTaskServices(workspaceId, filters, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch tasks",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(result)

}
