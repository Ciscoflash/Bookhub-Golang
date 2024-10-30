package router

import (
	"github.com/cisco-flash/user-management-system/controller"
	"github.com/cisco-flash/user-management-system/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	version := app.Group("/api/v1")
	auth := version.Group("/auth")
	book := version.Group("/book")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)
	// this profile downwards automatically uses the middle ware
	app.Use(middleware.Authentication)
	auth.Get("/profile", controller.Profile)
	book.Post("/create", controller.CreateBook)
	book.Get("/all", controller.GetBooks)
	book.Get("/:id", controller.GetBook)
	book.Patch("/:id", controller.EditBook)
	book.Delete("/:id", controller.Delete)
}
