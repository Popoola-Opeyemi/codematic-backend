package gateways

import (
	"context"
	"fmt"

	"codematic/internal/thirdparty/paystack"
)

type PaystackProvider struct {
	apiKey  string
	baseURL string
	client  *paystack.Client
}

func NewPaystackProvider(baseURL, apiKey string, client *paystack.Client) *PaystackProvider {
	return &PaystackProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
}

func (p *PaystackProvider) InitDeposit(ctx context.Context, email string,
	req DepositRequest) (*InitDepositResponse, error) {
	resp, err := p.client.InitializeTransaction(&paystack.InitializeTransactionRequest{
		Amount:   req.Amount.String(),
		Email:    email,
		Metadata: req.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("paystack init deposit error: %w", err)
	}

	return &InitDepositResponse{
		AuthorizationURL: resp.Data.AuthorizationURL,
		Reference:        resp.Data.Reference,
		Provider:         paystack.ProviderPaystack,
	}, nil
}

func (p *PaystackProvider) VerifyTransaction(ctx context.Context, reference string) (*VerifyResponse, error) {
	resp, err := p.client.VerifyTransaction(reference)
	if err != nil {
		return nil, fmt.Errorf("paystack verify error: %w", err)
	}

	status := "failed"
	if resp.Data.Status == "success" {
		status = "success"
	}

	return &VerifyResponse{
		Provider:  paystack.ProviderPaystack,
		Status:    status,
		Amount:    resp.Data.Amount,
		Currency:  resp.Data.Currency,
		Reference: resp.Data.Reference,
		Raw:       resp.Data,
	}, nil
}
