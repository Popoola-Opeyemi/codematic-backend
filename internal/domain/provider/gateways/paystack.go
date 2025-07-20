package gateways

import (
	"context"
	"fmt"

	"codematic/internal/thirdparty/paystack"

	"go.uber.org/zap"
)

type PaystackProvider struct {
	client *paystack.Client
	logger *zap.Logger
}

func NewPaystackProvider(logger *zap.Logger, baseURL, apiKey string) *PaystackProvider {
	return &PaystackProvider{
		client: paystack.NewPaystackClient(logger, baseURL, apiKey),
		logger: logger,
	}
}

func (p *PaystackProvider) InitDeposit(ctx context.Context,
	req DepositRequest) (GatewayResponse, error) {
	resp, err := p.client.InitializeTransaction(&paystack.InitializeTransactionRequest{
		Amount:   req.Amount.String(),
		Email:    req.Email,
		Metadata: req.Metadata,
	})
	if err != nil {
		return GatewayResponse{}, fmt.Errorf("paystack init deposit error: %w", err)
	}

	p.logger.Sugar().Infow("Paystack", "response", fmt.Sprintf("%+v", resp))

	return GatewayResponse{
		AuthorizationURL: resp.Data.AuthorizationURL,
		Reference:        resp.Data.Reference,
		Provider:         paystack.ProviderPaystack,
		ProviderID:       req.ProviderID,
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
