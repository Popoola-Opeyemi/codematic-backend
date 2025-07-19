package gateways

import (
	"codematic/internal/thirdparty/flutterwave"
)

type FlutterwaveProvider struct {
	apiKey  string
	baseURL string
	client  *flutterwave.Client
}

func NewFlutterwaveProvider(baseURL, apiKey string, client *flutterwave.Client) *FlutterwaveProvider {
	return &FlutterwaveProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
	}
}

// func (f *FlutterwaveProvider) InitDeposit(ctx context.Context, email string, req DepositRequest) (*InitDepositResponse, error) {
// 	txRef, ok := req.Metadata["reference"].(string)
// 	if !ok || txRef == "" {
// 		return nil, fmt.Errorf("missing reference in metadata")
// 	}

// 	amountFloat, err := strconv.ParseFloat(req.Amount.String(), 64)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid amount format: %w", err)
// 	}

// 	resp, err := f.client.InitializePayment(&flutterwave.InitPaymentRequest{
// 		TxRef:       txRef,
// 		Amount:      amountFloat,
// 		Currency:    req.Currency,
// 		RedirectURL: req.RedirectURL,
// 		Customer: map[string]string{
// 			"email": email,
// 		},
// 		Meta: req.Metadata,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("flutterwave init deposit error: %w", err)
// 	}

// 	return &InitDepositResponse{
// 		AuthorizationURL: resp.Data.Link,
// 		Reference:        txRef,
// 		Provider:         flutterwave.ProviderFlutterwave,
// 	}, nil
// }

// func (f *FlutterwaveProvider) VerifyTransaction(ctx context.Context, reference string) (*VerifyResponse, error) {
// 	txID, err := strconv.Atoi(reference)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid reference format: expected int txID")
// 	}

// 	resp, err := f.client.VerifyPayment(txID)
// 	if err != nil {
// 		return nil, fmt.Errorf("flutterwave verify error: %w", err)
// 	}

// 	status := "failed"
// 	if resp.Data.Status == "successful" {
// 		status = "success"
// 	}

// 	return &VerifyResponse{
// 		Provider:  flutterwave.ProviderFlutterwave,
// 		Status:    status,
// 		Amount:    int64(resp.Data.Amount),
// 		Currency:  resp.Data.Currency,
// 		Reference: resp.Data.TxRef,
// 		Raw:       resp.Data,
// 	}, nil
// }
