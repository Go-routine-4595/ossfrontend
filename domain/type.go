package domain

// Router defines the router
type Router struct {
	// in: query
	RouterID            string `json:"router-id" form:"router-id"`
	RouterSerial        string `json:"router-serial" form:"router-serial"`
	OperatorName        string `json:"operator-name" form:"operator-name"`
	IsoCountryCode      string `json:"iso-country-code" form:"iso-country-code"`
	Mac                 string `json:"mac" form:"mac"`
	RouterModel         string `json:"router-model" form:"router-model"`
	AccountID           string `json:"account-id" form:"account-id"`
	AgentLastConnection string `json:"agent-last-connection"`
	AgentVersion        string `json:"agent-version"`
}

// Account defines the account
type Account struct {
	AccountID   string `json:"account-id" form:"account-id"`
	Status      string `json:"status" form:"status"`
	AccountType string `json:"account-type" form:"account-type"`
}

// AccountDef defines the link between account and router
type AccountDef struct {
	Account
	routers []Router
}

// Command to be executed for an account
type Command struct {
	AccountID string `json:"account-id" form:"account-id"`
}

// Response to a Get Paged Router
type Response struct {
	Last    int      `json:"last"`
	Routers []Router `json:"routers"`
}
