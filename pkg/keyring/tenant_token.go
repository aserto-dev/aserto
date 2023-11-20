package keyring

import "time"

type TenantToken struct {
	ExpiresAt           time.Time `json:"expires_at"` // UTC timestamp when access_token expires
	TenantID            string    `json:"tenant_id"`
	AuthorizerAPIKey    string    `json:"authorizer_api_key"`
	RegistryDownloadKey string    `json:"registry_download_key"`
	RegistryUploadKey   string    `json:"registry_upload_key"`
	DecisionLogsKey     string    `json:"decision_logs_key"`
	DirectoryReadKey    string    `json:"directory_read_key"`
	DirectoryWriteKey   string    `json:"directory_write_key"`
	DiscoveryKey        string    `json:"discovery_key"`
}

func (t *TenantToken) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}
