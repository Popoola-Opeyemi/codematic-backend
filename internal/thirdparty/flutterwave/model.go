package flutterwave

type InitPaymentRequest struct {
	TxRef          string                 `json:"tx_ref"`
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	RedirectURL    string                 `json:"redirect_url"`
	PaymentOptions string                 `json:"payment_options,omitempty"`
	Customer       map[string]string      `json:"customer"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
}

type InitPaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

type VerifyPaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID            int     `json:"id"`
		TxRef         string  `json:"tx_ref"`
		FlwRef        string  `json:"flw_ref"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		ChargedAmount float64 `json:"charged_amount"`
		Status        string  `json:"status"`
	} `json:"data"`
}

type InitializeTransactionRequest struct {
	TxRef    string                 `json:"tx_ref"`
	Amount   string                 `json:"amount"`
	Currency string                 `json:"currency"`
	Customer Customer               `json:"customer"`
	Meta     map[string]interface{} `json:"meta"`
}

type Customer struct {
	Email string `json:"email"`
}

type InitializeTransactionResponse struct {
	Status string `json:"status"`
	Data   struct {
		Link  string `json:"link"`
		TxRef string `json:"tx_ref"`
	} `json:"data"`
}

type VerifyTransactionResponse struct {
	Status string `json:"status"`
	Data   struct {
		Status   string `json:"status"` // "successful"
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		TxRef    string `json:"tx_ref"`
	} `json:"data"`
}
