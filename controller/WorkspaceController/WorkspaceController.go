package workspacecontroller

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"rishik.com/services/workspace"
)

func WorkspaceController(c fiber.Ctx) error {
	userId := c.Locals("userId")
	var body workspace.WorkspaceRequest
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Name == "" || userId == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")

	}
	workspace, err := workspace.CreateWorkSpaceService(body, userId.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something Went wrong")

	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Workspace Created sucessfully",
		"workspace": workspace,
	})
}

type UpdateWorkspace struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func UpdateWorkSpacecontroller(c fiber.Ctx) error {
	worksapceId := c.Params("id")
	var body UpdateWorkspace
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	workspace, err := workspace.UpdateWorkspace(worksapceId, body.Name, body.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update workspace",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"workspace": workspace,
	})

}
func DeletWorkspace(c fiber.Ctx) error {
	workspaceidStr := c.Params("id")
	userId := c.Locals("userId")
	if workspaceidStr == "" || userId == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Missing workspace ID or user ID")
	}
	workspaceId, err := primitive.ObjectIDFromHex(workspaceidStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid workspace ID format")
	}
	userIdstr, ok := userId.(string)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format")
	}
	currentWorkspace, err := workspace.DeleteWorkspaceService(workspaceId.Hex(), userIdstr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	response := fiber.Map{
		"message":           "Workspace deleted successfully",
		"curreentWorkspace": currentWorkspace,
	}
	return c.Status(fiber.StatusOK).JSON(response)

}
func GetUserWorkspace(c fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	workspace, err := workspace.GetAllWorkSpacesUserIsMemberService(userId.(string))
	if err != nil {
		log.Printf("the error is %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workspaces",
		})
	}
	return c.JSON(fiber.Map{
		"wprkspaces": workspace,
	})
}
func GetUserWorkspacebyId(c fiber.Ctx) error {
	idstr := c.Params("id")
	if idstr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": " ID of workspace is required",
		})
	}
	workspace, err := workspace.GetWorkSpacebyIdService(idstr)
	if err != nil {
		log.Printf("the error is %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workspaces",
		})
	}
	return c.JSON(fiber.Map{
		"wprkspaces": workspace,
	})
}
func GetWorkspaceAnalytics(c fiber.Ctx) error {
	workspaceId := c.Params("workspaceId")
	if workspaceId == "" {
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
	analytics, err := workspace.GetWorkspaceAnalyticsServices(workspaceIdobj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workspaces",
		})

	}
	return c.JSON(analytics)

}
