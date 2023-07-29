package domain

// Router defines the router
type Router struct {
	// in: query
	RouterSerial string `json:"router-serial" form:"router-serial"`
}

// Account defines the account
type Account struct {
	AccountID string `json:"account-id" form:"account-id"`
}

// AccountDef defines the link between account and router
type AccountDef struct {
	Account
	routers []Router
}
