package events

type AlertTriggeredEvent struct {
	AlertID     string `json:"alert_id"`
	UserID      string `json:"user_id"`
	TokenSymbol string `json:"token_symbol"`
	Price       string `json:"price"`
	Chain       string `json:"chain"`
	Timestamp   int64  `json:"timestamp"`
}
