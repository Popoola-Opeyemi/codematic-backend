package gateways

import (
	"codematic/internal/thirdparty/flutterwave"
)

type FlutterwaveProvider struct {
	apiKey  string
	baseURL string
	client  *flutterwave.Client
}

func NewFlutterwaveProvider(baseURL, apiKey string,
	client *flutterwave.Client) *FlutterwaveProvider {
	return &FlutterwaveProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
}
