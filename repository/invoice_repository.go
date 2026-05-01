package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"invoice-view/model"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *model.Invoice) error
	GetByID(ctx context.Context, id string) (*model.Invoice, error)
}

type invoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *model.Invoice) error {
	customerJSON, err := json.Marshal(invoice.Customer)
	if err != nil {
		return fmt.Errorf("marshal customer: %w", err)
	}

	itemsJSON, err := json.Marshal(invoice.Items)
	if err != nil {
		return fmt.Errorf("marshal items: %w", err)
	}

	summaryJSON, err := json.Marshal(invoice.Summary)
	if err != nil {
		return fmt.Errorf("marshal summary: %w", err)
	}

	metadataJSON, err := json.Marshal(invoice.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	configJSON, err := json.Marshal(invoice.Config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	query := `
        INSERT INTO invoices 
            (id, invoice_number, currency, amount, status, issued_at, expired_at, customer, items, summary, metadata, config, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.db.ExecContext(ctx, query,
		invoice.ID,
		invoice.InvoiceNumber,
		invoice.Currency,
		invoice.Amount,
		invoice.Status,
		invoice.IssuedAt.Format(time.RFC3339),
		invoice.ExpiredAt.Format(time.RFC3339),
		string(customerJSON),
		string(itemsJSON),
		string(summaryJSON),
		string(metadataJSON),
		string(configJSON),
		invoice.CreatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return fmt.Errorf("insert invoice: %w", err)
	}

	return nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id string) (*model.Invoice, error) {
	query := `
        SELECT id, invoice_number, currency, amount, status, issued_at, expired_at, 
               customer, items, summary, metadata, config, created_at
        FROM invoices 
        WHERE id = ?`

	var row model.InvoiceRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&row.ID,
		&row.InvoiceNumber,
		&row.Currency,
		&row.Amount,
		&row.Status,
		&row.IssuedAt,
		&row.ExpiredAt,
		&row.Customer,
		&row.Items,
		&row.Summary,
		&row.Metadata,
		&row.Config,
		&row.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("scan invoice: %w", err)
	}

	// Parse waktu
	issuedAt, _ := time.Parse(time.RFC3339, row.IssuedAt)
	expiredAt, _ := time.Parse(time.RFC3339, row.ExpiredAt)
	createdAt, _ := time.Parse(time.RFC3339, row.CreatedAt)

	// Parse JSON fields
	var customer model.Customer
	_ = json.Unmarshal([]byte(row.Customer), &customer)

	var items []model.Item
	_ = json.Unmarshal([]byte(row.Items), &items)

	var summary model.Summary
	_ = json.Unmarshal([]byte(row.Summary), &summary)

	var metadata model.Metadata
	_ = json.Unmarshal([]byte(row.Metadata), &metadata)

	var config model.Config
	_ = json.Unmarshal([]byte(row.Config), &config)

	invoice := &model.Invoice{
		InvoiceRequest: model.InvoiceRequest{
			ID:            row.ID,
			InvoiceNumber: row.InvoiceNumber,
			Currency:      row.Currency,
			Amount:        row.Amount,
			Status:        row.Status,
			IssuedAt:      issuedAt,
			ExpiredAt:     expiredAt,
			Customer:      customer,
			Items:         items,
			Summary:       summary,
			Metadata:      metadata,
			Config:        config,
		},
		CreatedAt: createdAt,
	}

	return invoice, nil
}
