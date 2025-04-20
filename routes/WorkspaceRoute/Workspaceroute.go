package workspaceroute

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/asynchandler"
	workspacecontroller "rishik.com/controller/WorkspaceController"

	"rishik.com/middleware"
)

func SetupWorkspaceRoutes(app fiber.Router) {
	workspace := app.Group("/workspace")
	workspace.Post("/create", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.WorkspaceController))
	workspace.Patch("/:id", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.UpdateWorkSpacecontroller))
	workspace.Delete("/:id", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.DeletWorkspace))
	workspace.Get("/workspaces", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.GetUserWorkspace))
	workspace.Get("/:id", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.GetUserWorkspacebyId))
	workspace.Get("/workspaceAnalytics/:workspaceId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(workspacecontroller.GetWorkspaceAnalytics))
}
