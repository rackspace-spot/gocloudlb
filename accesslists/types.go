package accesslists

type CreateOpts struct {
	Address string `json:"address"`
	// The type of the access list, e.g., "ALLOW" or "DENY"
	Type string `json:"type"`
}

type NetworkItem struct {
	Address string `json:"address"`
	// The type of the access list, e.g., "ALLOW" or "DENY"
	Type string `json:"type"`
	ID   uint64 `json:"id"`
}
