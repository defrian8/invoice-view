package model

import "time"

type InvoiceRequest struct {
	ID            string    `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	Currency      string    `json:"currency"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiredAt     time.Time `json:"expired_at"`
	Customer      Customer  `json:"customer"`
	Items         []Item    `json:"items"`
	Summary       Summary   `json:"summary"`
	Metadata      Metadata  `json:"metadata"`
	Config        Config    `json:"config"`
}

type Customer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type Item struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
	Total float64 `json:"total"`
}

type Summary struct {
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Discount   float64 `json:"discount"`
	GrandTotal float64 `json:"grand_total"`
}

type Metadata struct {
	Notes string `json:"notes"`
}

type Config struct {
	Template    string `json:"template"`
	Orientation string `json:"orientation"`
}

type Invoice struct {
	InvoiceRequest
	CreatedAt time.Time `json:"created_at"`
}

type InvoiceRow struct {
	ID            string
	InvoiceNumber string
	Currency      string
	Amount        float64
	Status        string
	IssuedAt      string
	ExpiredAt     string
	Customer      string
	Items         string
	Summary       string
	Metadata      string
	Config        string
	CreatedAt     string
}
