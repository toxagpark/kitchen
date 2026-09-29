package restaurant_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

type Client struct {
	orderURL   string
	httpClient *http.Client
}

func NewClient(
	cfg Config,
	httpClient *http.Client,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		orderURL:   cfg.BaseURL + cfg.OrderPath,
		httpClient: httpClient,
	}
}

type createOrderRequest struct {
	OrderID    uuid.UUID          `json:"order_id"`
	Items      []orderItemRequest `json:"items"`
	TotalPrice int64              `json:"total_price"`
}

type orderItemRequest struct {
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unit_price"`
	TotalPrice int64  `json:"total_price"`
}

func (c *Client) CreateOrder(
	ctx context.Context,
	order domain.Order,
	items []domain.OrderItem,
) error {
	requestBody := createOrderRequest{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
		Items:      make([]orderItemRequest, 0, len(items)),
	}

	for _, item := range items {
		requestBody.Items = append(requestBody.Items, orderItemRequest{
			Name:       item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		})
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal restaurant order request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.orderURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create restaurant order request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send restaurant order request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"restaurant returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
