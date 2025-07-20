package flutterwave

import (
	"codematic/internal/thirdparty/baseclient"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	client  baseclient.BaseClient
	secret  string
	baseURL string
	logger  *zap.Logger
}

func NewFlutterwaveClient(baseURL, secret string, logger *zap.Logger) *Client {
	config := baseclient.DefaultConfig()
	config.Timeout = 60 * time.Second
	config.MaxIdleConns = 50
	config.MaxIdleConnsPerHost = 5

	baseCLI := baseclient.NewBaseClientWithConfig(logger, config)
	return &Client{client: *baseCLI, secret: secret, baseURL: baseURL, logger: logger}
}

// Request & response structs

func (c *Client) InitializePayment(req *InitPaymentRequest) (*InitPaymentResponse, error) {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	resp, err := c.client.MakeRequest("POST", url, req, c.authHeaders())
	if err != nil {
		c.logger.Error("init payment request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read init response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("init payment error: %s", string(body))
	}

	var out InitPaymentResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("unmarshal init response: %w", err)
	}

	return &out, nil
}

func (c *Client) VerifyPayment(txID int) (*VerifyPaymentResponse, error) {
	url := fmt.Sprintf("%s/transactions/%d/verify", c.baseURL, txID)

	resp, err := c.client.MakeRequest("GET", url, nil, c.authHeaders())
	if err != nil {
		c.logger.Error("verify payment request failed", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read verify response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verify payment error: %s", string(body))
	}

	var out VerifyPaymentResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("unmarshal verify response: %w", err)
	}

	return &out, nil
}

func (c *Client) authHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.secret,
		"Content-Type":  "application/json",
	}
}
