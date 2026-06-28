package main

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	"time"

	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"santiirepair.dev/git-receipt/cache"
	"santiirepair.dev/git-receipt/github"
	"santiirepair.dev/git-receipt/handlers"
)

var (
	servers = []string{"Grace Hopper", "Alan Turing", "Ada Lovelace", "Tim Berners-Lee", "Linus Torvalds"}
)

//go:embed template.svg
var templateFS embed.FS

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	templateBytes, err := templateFS.ReadFile("template.svg")
	if err != nil {
		log.Printf("❌ Error loading template: %v\n", err)
		os.Exit(1)
	}
	templateContent := string(templateBytes)

	githubService := github.NewService()
	cacheManager := cache.NewGitHubCacheManager(githubService)
	receiptHandler := handlers.NewReceiptHandler(cacheManager, templateContent, servers)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		AppName:               "GitHub Receipt",
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		JSONEncoder: func(v any) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	handlers.SetupRoutes(app, receiptHandler, cacheManager, githubService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
