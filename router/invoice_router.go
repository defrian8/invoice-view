package router

import (
	"github.com/gofiber/fiber/v2"

	"invoice-view/model"
	"invoice-view/router/middleware"
	"invoice-view/service"
)

type invoiceHandler struct {
	service         service.InvoiceService
	templateService service.TemplateService
}

func NewInvoiceRouter(app *fiber.App, invoiceSvc service.InvoiceService, templateScv service.TemplateService) {
	h := &invoiceHandler{
		service:         invoiceSvc,
		templateService: templateScv,
	}

	app.Get("/invoice/:id", middleware.RateLimiter(), h.GetInvoiceHTML)

	v1 := app.Group("/v1")
	v1.Post("/invoice", middleware.APIKeyMiddleware(), h.CreateInvoice)
}

func (h *invoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var req model.InvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	invoice, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

func (h *invoiceHandler) GetInvoiceHTML(c *fiber.Ctx) error {
	id := c.Params("id")

	invoice, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "invoice not found" {
			return c.Status(fiber.StatusNotFound).SendString("Invoice not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	buf, err := h.templateService.Render(c.Context(), invoice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to render template")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(buf.Bytes())
}
