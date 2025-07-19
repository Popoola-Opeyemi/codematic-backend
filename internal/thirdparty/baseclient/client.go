package baseclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BaseClient struct {
	HTTPClient *http.Client
}

type Config struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Timeout:             30 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}

func NewBaseClient(timeout time.Duration) *BaseClient {
	config := DefaultConfig()
	config.Timeout = timeout
	return NewBaseClientWithConfig(config)
}

func NewBaseClientWithConfig(config Config) *BaseClient {
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,
		DisableCompression:  false,
	}

	return &BaseClient{
		HTTPClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
	}
}

func (bc *BaseClient) GetHTTPClient() *http.Client {
	return bc.HTTPClient
}

func (bc *BaseClient) MakeRequest(method, url string, body interface{},
	headers map[string]string) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return bc.HTTPClient.Do(req)
}
