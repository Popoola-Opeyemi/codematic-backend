package idempotency

type (
	CreateParams struct {
		ID           string
		TenantID     string
		UserID       string
		Key          string
		Endpoint     string
		RequestHash  string
		ResponseBody map[string]interface{}
		StatusCode   int
	}

	UpdateResponseParams struct {
		TenantID     string
		Key          string
		Endpoint     string
		ResponseBody map[string]interface{}
		StatusCode   int
	}
)
