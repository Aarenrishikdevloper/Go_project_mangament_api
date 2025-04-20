package taskroute

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/asynchandler"
	"rishik.com/controller/taskcontroller"

	"rishik.com/middleware"
)

func SetupTaskRoutes(app fiber.Router) {
	task := app.Group("/task")

	task.Post("/workspace/:workspaceId/projects/:projectId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(taskcontroller.CreateTask))
	task.Patch("/workspace/:workspaceId/projects/:projectId/task/:taskId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(taskcontroller.UpdateTaskController))
	task.Delete("/workspace/:workspaceId/task/:taskId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(taskcontroller.DeleteTaskController))
	task.Get("/workspace/:workspaceId/projects/:projectId/task/:taskId", middleware.IsAuthenticated(), asynchandler.AsyncHandler(taskcontroller.GetTaskBYIdController))
	task.Get("/workspaces/:workspaceId/task", middleware.IsAuthenticated(), asynchandler.AsyncHandler(taskcontroller.GetAllWorkspace))

}
