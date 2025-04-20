package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"

	"rishik.com/asynchandler" // Assuming this is the correct import path for your async handler package
	"rishik.com/db"
	"rishik.com/middleware"
	workspaceroute "rishik.com/routes/WorkspaceRoute"
	"rishik.com/routes/authroute"
	"rishik.com/routes/memberroute"
	"rishik.com/routes/projectroute"
	"rishik.com/routes/taskroute"
	"rishik.com/routes/userroute"
	// Assuming this is the correct import path for your auth routes package
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("warning:Couldn't load .env file", err)
	}
	uri := os.Getenv("MONGO_URI")
	err = db.InitDB(uri, "teamsync_db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true // allow all origins dynamically
		},
		AllowCredentials: true,
	}))

	api := app.Group("/api")
	authroute.SetupAuthRoutes(api)
	app.Use(middleware.IsAuthenticated())
	userroute.SetupUserRoutes(api)
	workspaceroute.SetupWorkspaceRoutes(api)
	projectroute.SetupProjectRoutes(api)
	taskroute.SetupTaskRoutes(api)
	memberroute.SetupMemberRoutes(api)
	app.Get("/", asynchandler.AsyncHandler(func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	}))

	log.Fatal(app.Listen(":8000"))
}
