package paystack

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
	apiKey  string
	logger  *zap.Logger
	baseURL string
	baseclient.BaseClient
}

func NewPaystackClient(logger *zap.Logger, baseURL, apiKey string) *Client {
	config := baseclient.DefaultConfig()
	config.Timeout = 60 * time.Second
	config.MaxIdleConns = 50
	config.MaxIdleConnsPerHost = 5

	baseClient := baseclient.NewBaseClientWithConfig(logger, config)

	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		BaseClient: *baseClient,
		logger:     logger, // ✅ ensure logger is assigned
	}
}

func (c *Client) InitializeTransaction(req *InitializeTransactionRequest) (*InitializeTransactionResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", c.baseURL)

	resp, err := c.MakeRequest(http.MethodPost, url, req, c.authHeaders())
	if err != nil {
		c.logger.Error("initialize transaction request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("initialize transaction failed: %s", string(bodyBytes))
	}

	var initResp InitializeTransactionResponse
	if err := json.Unmarshal(bodyBytes, &initResp); err != nil {
		return nil, fmt.Errorf("unmarshal initialize response failed: %w", err)
	}

	return &initResp, nil
}

func (c *Client) VerifyTransaction(reference string) (*VerifyTransactionResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", c.baseURL, reference)

	resp, err := c.MakeRequest(http.MethodGet, url, nil, c.authHeaders())
	if err != nil {
		c.logger.Error("verify transaction request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verify transaction failed: %s", string(bodyBytes))
	}

	var verifyResp VerifyTransactionResponse
	if err := json.Unmarshal(bodyBytes, &verifyResp); err != nil {
		return nil, fmt.Errorf("unmarshal verify response failed: %w", err)
	}

	return &verifyResp, nil
}

func (c *Client) authHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}
}

func readResponseBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(body)
}
