package idempotency

import "encoding/json"

type (
	CreateParams struct {
		ID           string
		TenantID     string
		UserID       string
		Key          string
		Endpoint     string
		RequestHash  string
		ResponseBody json.RawMessage

		StatusCode int
	}

	UpdateResponseParams struct {
		TenantID     string
		Key          string
		Endpoint     string
		ResponseBody json.RawMessage

		StatusCode int
	}
)
