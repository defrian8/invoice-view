package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"invoice-view/model"
	"invoice-view/repository"
)

type InvoiceService interface {
	Create(ctx context.Context, req *model.InvoiceRequest) (*model.Invoice, error)
	GetByID(ctx context.Context, id string) (*model.Invoice, error)
}

type invoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) Create(ctx context.Context, req *model.InvoiceRequest) (*model.Invoice, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("id is required")
	}
	if req.InvoiceNumber == "" {
		return nil, fmt.Errorf("invoice_number is required")
	}

	invoice := &model.Invoice{
		InvoiceRequest: *req,
		CreatedAt:      time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("service create: %w", err)
	}

	return invoice, nil
}

func (s *invoiceService) GetByID(ctx context.Context, id string) (*model.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("service get: %w", err)
	}

	return invoice, nil
}
