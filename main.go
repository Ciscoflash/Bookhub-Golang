package main

import (
	"log"

	"github.com/cisco-flash/user-management-system/database"
	"github.com/cisco-flash/user-management-system/router"
	"github.com/gofiber/fiber/v2"
	// "github.com/joho/godotenv"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2/middleware/cors"

)

var DB *gorm.DB

type User struct {
	firstName string `json:"firstname"`
	lastName  string `json:"lastname"`
	email     string `json:"email"`
	password  string `json:"password"`
}

func main() {
	// err := godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatal("error loading env File", err)
	// }
	database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
    	// Allows all request headers
   		 AllowHeaders: "content-type",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to USM Developer Api")
	})
	router.SetupRoutes(app)
	// SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
