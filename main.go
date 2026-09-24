package main

import (
	"fmt"
	"invoice-view/database"
	"invoice-view/repository"
	"invoice-view/router"
	"invoice-view/service"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	time.Local = loc

	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		panic("DB_HOST is required")
	}

	db := database.NewSQLite(fmt.Sprintf("file:%s?cache=shared&_journal_mode=WAL", dbHost))
	defer db.Close()

	invoiceRepo := repository.NewInvoiceRepository(db)

	invoiceSvc := service.NewInvoiceService(invoiceRepo)
	templateSvc := service.NewTemplateService()

	app := fiber.New(fiber.Config{
		AppName: "Invoice Service",
	})

	router.NewInvoiceRouter(app, invoiceSvc, templateSvc)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("[server] starting on :3000")
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
