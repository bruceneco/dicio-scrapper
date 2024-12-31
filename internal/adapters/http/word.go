package http

import (
	"dicio-scrapper/internal/domain/word"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Word struct {
	service word.Servicer
	app     *fiber.App
}

func NewWordController(service word.Servicer, app *fiber.App) {
	c := Word{service, app}

	group := app.Group("/word")
	group.Get("/most-searched", c.MostSearched)
	group.Get("", c.GetByContent)
}

func (w *Word) MostSearched(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "invalid page",
		})
	}
	numWords, err := w.service.EnqueueMostSearched(c.Context(), page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": fmt.Sprintf("failed to enqueue most searched of page %d", page),
		})
	}
	return c.JSON(fiber.Map{
		"words": numWords,
	})
}

func (w *Word) GetByContent(c *fiber.Ctx) error {
	content := c.Query("content")
	fmt.Println(content)
	if content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "content is empty",
			"message": "you must provide a word to be searched",
		})
	}
	wordResult, err := w.service.GetWord(c.Context(), content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": fmt.Sprintf("failed to get the word %s", content),
		})
	}
	return c.JSON(wordResult)
}
