package api

import "time"

type Token struct {
	Type                string    `json:"token_type"`
	Scope               string    `json:"scope"`
	Identity            string    `json:"id_token"`
	Access              string    `json:"access_token"`
	ExpiresIn           int       `json:"expires_in"`
	ExpiresAt           time.Time `json:"expires_at"` // UTC timestamp when access_token expires
	TenantID            string    `json:"tenant_id"`
	AuthorizerAPIKey    string    `json:"authorizer_api_key"`
	RegistryDownloadKey string    `json:"registry_download_key"`
	RegistryUploadKey   string    `json:"registry_upload_key"`
	DecisionLogsKey     string    `json:"decision_logs_key"`
}
