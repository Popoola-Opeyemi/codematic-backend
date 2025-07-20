package gateways

import (
	"context"
	"fmt"

	"codematic/internal/thirdparty/paystack"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type PaystackProvider struct {
	client    *paystack.Client
	logger    *zap.Logger
	apiSecret string
}

func NewPaystackProvider(logger *zap.Logger, baseURL, apiSecret string) *PaystackProvider {
	return &PaystackProvider{
		client:    paystack.NewPaystackClient(logger, baseURL, apiSecret),
		logger:    logger,
		apiSecret: apiSecret,
	}
}

func (p *PaystackProvider) InitDeposit(ctx context.Context,
	req DepositRequest) (GatewayResponse, error) {
	amountInKobo := req.Amount.Mul(decimal.NewFromInt(100))

	resp, err := p.client.InitializeTransaction(&paystack.InitializeTransactionRequest{
		Amount:   amountInKobo.String(),
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

func (p *PaystackProvider) VerifyTransaction(ctx context.Context,
	reference string) (*VerifyResponse, error) {
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

func (p *PaystackProvider) VerifyWebhookSignature(body []byte,
	signatureHeader string) (bool, error) {

	expectedSig, err := paystack.GenerateSignature(p.apiSecret, body)
	if err != nil {
		return false, fmt.Errorf("failed to generate signature: %w", err)
	}

	p.logger.Sugar().Debug("signatureHeader ", signatureHeader)
	p.logger.Sugar().Debug("expectedSig ", expectedSig)

	return expectedSig == signatureHeader, nil
}
