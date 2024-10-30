package controller

import (
	"fmt"
	"net/http"

	"github.com/cisco-flash/user-management-system/database"
	"github.com/cisco-flash/user-management-system/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Books struct {
	Title         string
	Description   string
	Author        string
	Publisher     string
	CoverImage    string
	BookUrl       string
	Category      string
	Rating        int
	PublisherName string
}

func CreateBook(context *fiber.Ctx) error {
	var req Books

	if err := context.BodyParser(&req); err != nil {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}
	token := context.Locals("jwt").(*jwt.Token)

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to extract claims"})
	}

	userId := claims["user_id"]
	userIdStr := fmt.Sprintf("%.0f", userId)

	book := models.Books{
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		Publisher:   userIdStr,
		CoverImage:  req.CoverImage,
		Category:    req.Category,
		BookUrl:     req.BookUrl,
		Rating:      req.Rating,
	}

	res := database.DB.Create(&book)
	if res.Error != nil {
		return context.Status(400).SendString(res.Error.Error())
	}

	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "book added"})
}

func GetBooks(context *fiber.Ctx) error {
	var books []models.Books

	// Use the Find method and check the error returned by the `Error` field
	result := database.DB.
		Table("books").
		Select("books.*, users.first_name || ' ' || users.last_name AS publisher_name").
		Joins("JOIN users ON users.id = books.publisher::int").
		Find(&books)

	if result.Error != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": result.Error.Error()})
	}
	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "Data retrieved", "data": books})
}

func GetBook(context *fiber.Ctx) error {
	id := context.Params("id")
	var book models.Books

	if id == "" {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "id is missing"})
	}

	res := database.DB.
		Table("books").
		Select("books.*, CONCAT(users.first_name, ' ', users.last_name) AS publisher_name").
		Joins("JOIN users ON users.id = books.publisher::int").
		First(&book, id).Error
	if res != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid id provided"})
	}
	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "data retrieved", "data": book})
}

func EditBook(context *fiber.Ctx) error {
	id := context.Params("id")
	var book Books
	if err := context.BodyParser(&book); err != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if id == "" {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "id is missong"})
	}

	var oldBook models.Books
	database.DB.Find(&oldBook, id)
	if oldBook.Title == "" {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "no book found with this id"})
	}
	res := database.DB.Model(&oldBook).Updates(book).Error
	if res != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error updating book"})
	}
	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "book updated"})
}

func Delete(context *fiber.Ctx) error {
	id := context.Params("id")
	var book models.Books
	if id == "" {
		return context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "id is missing"})
	}
	err := database.DB.Find(&book, id).Error
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "no book found"})
	}
	database.DB.Delete(&book, id)
	return context.Status(http.StatusOK).JSON(fiber.Map{"message": "book deleted"})
}
