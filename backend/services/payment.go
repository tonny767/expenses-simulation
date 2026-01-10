package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend/constants"
	"backend/models"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, expense *models.Expense, approval *models.Approval) (*models.Expense, *models.Approval, error)
}

type paymentService struct {
	client  *http.Client
	baseURL string
}

func NewPaymentService() PaymentService {
	return &paymentService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: os.Getenv("PAYMENT_BASE_URL"),
	}
}

// ONLY FOR TESTING PURPOSES
func NewPaymentServiceWithBaseURL(baseURL string) PaymentService {
	return &paymentService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
	}
}

type PaymentRequest struct {
	Amount     int64  `json:"amount"`
	ExternalID string `json:"external_id"`
}

type PaymentResponse struct {
	Data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

// ProcessPayment calls the external processor and updates the structs in-memory
func (s *paymentService) ProcessPayment(ctx context.Context, expense *models.Expense, approval *models.Approval) (*models.Expense, *models.Approval, error) {
	if s.baseURL == "" {
		return nil, nil, errors.New("PAYMENT_BASE_URL not configured")
	}

	reqBody := PaymentRequest{
		Amount:     expense.AmountIDR,
		ExternalID: expense.UUID.String(),
	}

	payload, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/v1/payments", bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("payment processor request failed: %w", err)
	}
	defer resp.Body.Close()

	var result PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode payment response: %w", err)
	}

	// Only proceed if payment succeeded or idempotency
	if !(resp.StatusCode == http.StatusOK || (resp.StatusCode == http.StatusBadRequest && result.Message == "external id already exists")) {
		return nil, nil, fmt.Errorf("payment failed with status %d: %s", resp.StatusCode, result.Message)
	}

	now := time.Now()

	if expense.AutoApproved {
		expense.Status = constants.ExpenseStatusCompleted
		approval.Notes = "auto-approved"
	} else {
		expense.Status = constants.ExpenseStatusCompleted
	}

	expense.ProcessedAt = &now

	return expense, approval, nil
}
