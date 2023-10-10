package api

import "time"

type Token struct {
	Type            string    `json:"token_type"`
	Scope           string    `json:"scope"`
	Identity        string    `json:"id_token"`
	Access          string    `json:"access_token"`
	Subject         string    `json:"subject"`
	ExpiresIn       int       `json:"expires_in"`
	ExpiresAt       time.Time `json:"expires_at"` // UTC timestamp when access_token expires
	DefaultTenantID string    `json:"default_tenant_id"`
}

func (t *Token) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}
