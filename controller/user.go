package controller

import (
	"net/http"

	"github.com/cisco-flash/user-management-system/database"
	"github.com/cisco-flash/user-management-system/models"
	"github.com/cisco-flash/user-management-system/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type authRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(context *fiber.Ctx) error {
	var req authRequest

	if err := context.BodyParser(&req); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}

	// Populate the user struct with data from the request
	user := models.User{
		FirstName: req.FirstName, // Assuming your authRequest has these fields
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  utils.HashPassword(req.Password),
	}

	// Ensure the database is initialized
	if database.DB == nil {
		return context.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Database not initialized"})
	}

	res := database.DB.Create(&user)
	if res.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": res.Error.Error()})
	}
	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "User created"})
}

func Login(context *fiber.Ctx) error {
	var req loginRequest

	if err := context.BodyParser(&req); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}
	var user models.User

	if res := database.DB.Where("email = ?", req.Email).First(&user); res.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "user not found"})
	}
	// fmt.Println("res", user)
	if !utils.ComparePassword(user.Password, req.Password) {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "incorrect paasword"})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return context.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return context.JSON(fiber.Map{"message": "login success", "tokem": token})
}

func Profile(context *fiber.Ctx) error {
	token := context.Locals("jwt").(*jwt.Token)

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to extract claims"})
	}
	var user models.User
	userId := claims["user_id"]
	user.Password = ""
	if err := database.DB.First(&user, userId).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
	}

	return context.JSON(fiber.Map{"message": "profile", "data": user})
}
