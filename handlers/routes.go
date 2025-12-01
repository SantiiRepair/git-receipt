package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"santiirepair.dev/git-receipt/cache"
	"santiirepair.dev/git-receipt/github"
)

func SetupRoutes(app *fiber.App, receiptHandler *ReceiptHandler, cacheManager *cache.GitHubCacheManager) {
	githubService := github.NewGitHubService()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"message":           "Not Found",
			"documentation_url": "https://github.com/SantiiRepair/git-receipt",
		})
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		_, err := githubService.CheckAPIStatus()
		githubStatus := "ok"
		if err != nil {
			githubStatus = "error"
		}

		return c.JSON(fiber.Map{
			"status":      "ok",
			"timestamp":   time.Now().Unix(),
			"github_api":  githubStatus,
			"cache_stats": cacheManager.GetCacheMetrics(),
		})
	})

	app.Get("/cache/metrics", func(c *fiber.Ctx) error {
		return c.JSON(cacheManager.GetCacheMetrics())
	})

	app.Get("/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")

		if username == "" {
			return c.Status(400).SendString("Username is required")
		}

		html, err := receiptHandler.GenerateReceipt(username)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("User '%s' not found: %v", username, err))
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=300")
		c.Set("X-Cache-Status", "HIT")

		return c.SendString(html)
	})
}
