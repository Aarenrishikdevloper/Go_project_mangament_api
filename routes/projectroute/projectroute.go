package projectroute

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/asynchandler"
	projectcontroller "rishik.com/controller/projectController"
	"rishik.com/middleware"
)

func SetupProjectRoutes(app fiber.Router) {
	project := app.Group("/project")
	project.Post("/:id", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.CreateProject))
	project.Put("/workspace/:workspaceId/projects/:projectId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.UpdateProject))
	project.Delete("/workspace/:workspaceId/projects/:projectId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.DeleteProject))
	project.Get("/workspace/:workspaceId/projects/:projectId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.GetProjectBYidAndWorkspaceIdController))
	project.Get("/projects", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.GetProjectsinWorkspaces))
	project.Get("/projectAnalysis/project/:projectId/workspace/:workspaceId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(projectcontroller.GetProjectAnalytics))

}
