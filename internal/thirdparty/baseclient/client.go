package baseclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type BaseClient struct {
	HTTPClient *http.Client
	Logger     *zap.Logger
}

// Config holds configuration for the underlying HTTP client
type Config struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
}

// DefaultConfig returns a sensible default HTTP client config
func DefaultConfig() Config {
	return Config{
		Timeout:             30 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}

// NewBaseClient creates a new client with a custom timeout
func NewBaseClient(logger *zap.Logger, timeout time.Duration) *BaseClient {
	cfg := DefaultConfig()
	cfg.Timeout = timeout
	return NewBaseClientWithConfig(logger, cfg)
}

// NewBaseClientWithConfig creates a new client with full configuration
func NewBaseClientWithConfig(logger *zap.Logger, config Config) *BaseClient {
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,
		DisableCompression:  false,
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	return &BaseClient{
		Logger:     logger,
		HTTPClient: httpClient,
	}
}

// GetHTTPClient exposes the internal HTTP client
func (bc *BaseClient) GetHTTPClient() *http.Client {
	return bc.HTTPClient
}
func (bc *BaseClient) MakeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody *bytes.Buffer
	var serializedBody string

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		serializedBody = string(jsonData)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = &bytes.Buffer{}
	}

	// Prepare log fields
	logFields := []zap.Field{
		zap.String("method", method),
		zap.String("url", url),
		zap.String("body", serializedBody),
	}

	// Log headers individually
	for key, value := range headers {
		logFields = append(logFields, zap.String(fmt.Sprintf("header:%s", key), value))
	}

	bc.Logger.Info("Making HTTP request", logFields...)

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request failed: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		bc.Logger.Error("HTTP request failed", zap.Error(err))
		return nil, err
	}

	return resp, nil
}
