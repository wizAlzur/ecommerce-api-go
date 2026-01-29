package services

import (
	"bytes"
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, order *models.Order) (string, error)
}

type paymentService struct {
	cfg config.YooKassaConfig
}

func NewPaymentService(cfg config.YooKassaConfig) PaymentService {
	return &paymentService{cfg: cfg}
}

func (s *paymentService) CreatePayment(ctx context.Context, order *models.Order) (string, error) {
	payload := map[string]interface{}{
		"amount": map[string]interface{}{
			"value":    fmt.Sprintf("%.2f", order.TotalAmount),
			"currency": "RUB",
		},
		"confirmation": map[string]interface{}{
			"type":       "redirect",
			"return_url": s.cfg.SuccessURL,
		},
		"capture":     true,
		"description": fmt.Sprintf("Заказ №%d", order.ID),
		"metadata": map[string]interface{}{
			"order_id": fmt.Sprintf("%d", order.ID),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.cfg.ShopID, s.cfg.SecretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", fmt.Sprintf("%d-%d", order.ID, time.Now().UnixNano()))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("yookassa error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		ID           string `json:"id"`
		Confirmation struct {
			Type            string `json:"type"`
			ConfirmationURL string `json:"confirmation_url"`
		} `json:"confirmation"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("json decode failed: %w", err)
	}

	if result.Confirmation.Type != "redirect" {
		return "", fmt.Errorf("unexpected confirmation type: %s", result.Confirmation.Type)
	}

	return result.Confirmation.ConfirmationURL, nil
}
