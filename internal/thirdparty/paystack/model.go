package paystack

const (
	ProviderPaystack = "paystack"
)

type InitializeTransactionRequest struct {
	Email    string                 `json:"email"`
	Amount   string                 `json:"amount"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type InitializeTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type VerifyTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		Reference string `json:"reference"`
		// Add other fields as needed
	} `json:"data"`
}
