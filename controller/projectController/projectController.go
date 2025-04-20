package projectcontroller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"rishik.com/services/project"
)

func CreateProject(c fiber.Ctx) error {
	userId := c.Locals("userId")
	workspaceId := c.Params("id")
	var body project.ProjectInput
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Name == "" || userId == nil || workspaceId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")

	}
	project, err := project.CreateProject(userId.(string), workspaceId, body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something Went wrong")

	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Project Created sucessfully",
		"project": project,
	})
}
func UpdateProject(c fiber.Ctx) error {
	projectId := c.Params("projectId")
	workspaceId := c.Params("workspaceId")
	var body project.UpdateProjectBody
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Name == "" || workspaceId == "" || projectId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")

	}
	project, err := project.UpdatePeoject(workspaceId, projectId, body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something Went wrong")

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project updated sucessfully",
		"project": project,
	})
}
func DeleteProject(c fiber.Ctx) error {
	projectId := c.Params("projectId")
	workspaceId := c.Params("workspaceId")

	if projectId == "" || workspaceId == "" {

		return fiber.NewError(fiber.StatusBadRequest, "Missing workspace ID or user ID")
	}
	project, err := project.DeleteprojectService(workspaceId, projectId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something Went wrong")

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project Deleted sucessfully",
		"project": project,
	})
}
func GetProjectBYidAndWorkspaceIdController(c fiber.Ctx) error {
	projectId := c.Params("projectId")
	workspaceId := c.Params("workspaceId")

	if projectId == "" || workspaceId == "" {

		return fiber.NewError(fiber.StatusBadRequest, "Missing workspace ID or user ID")
	}
	project, err := project.GetProjectBYidAndWorkspaceIdService(workspaceId, projectId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something Went wrong")

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project Fetch Sucessfully",
		"project": project,
	})
}
func GetProjectsinWorkspaces(c fiber.Ctx) error {
	workspaceId := c.Query("workspaceId")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing workspace ID",
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
	result, err := project.GetProjectsInWorkspaces(workspaceId, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}
func GetProjectAnalytics(c fiber.Ctx) error {
	workspaceId := c.Params("workspaceId")
	projectId := c.Params("projectId")
	if workspaceId == "" || projectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": " ID of workspace is required",
		})

	}
	workspaceIdobj, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": " ID of workspace is Invalid",
		})
	}
	projectIdobj, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": " ID of workspace is Invalid",
		})
	}
	analytics, err := project.GetProjectAnalytics(workspaceIdobj, projectIdobj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workspaces",
		})

	}
	return c.JSON(analytics)

}
