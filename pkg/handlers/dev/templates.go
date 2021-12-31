package dev

type templateParams struct {
	TenantID     string
	PolicyName   string
	PolicyID     string
	DiscoveryURL string
	TenantKey    string
}

const configTemplate = templatePreamble + `
opa:
  instance_id: {{ .TenantID }}
  graceful_shutdown_period_seconds: 2
  local_bundles:
    paths: []
    skip_verification: true
  config:
    services:
      aserto-tenant:
        url: {{ .DiscoveryURL }}
        credentials:
          bearer:
            token: "{{ .TenantKey }}"
            scheme: "basic"
        headers:
          Aserto-Tenant-Id: {{ .TenantID }}
    discovery:
      name: "opa/discovery"
      prefix: "{{ .PolicyID }}"
`

const configTemplateLocal = templatePreamble + `
opa:
  instance_id: {{ .TenantID }}
  graceful_shutdown_period_seconds: 2
  local_bundles:
    paths: []
    skip_verification: true
`

const templatePreamble = `---
logging:
  prod: true
  log_level: info

directory_service:
  path: /app/db/directory.db

api:
  grpc:
    connection_timeout_seconds: 2
    listen_address: "0.0.0.0:8282"
    certs:
      tls_key_path: "/certs/grpc.key"
      tls_cert_path: "/certs/grpc.crt"
      tls_ca_cert_path: "/certs/grpc-ca.crt"
  gateway:
    listen_address: "0.0.0.0:8383"
    allowed_origins:
    - https://*.aserto.com
    - https://*aserto-console.netlify.app
    certs:
      tls_key_path: "/certs/gateway.key"
      tls_cert_path: "/certs/gateway.crt"
      tls_ca_cert_path: "/certs/gateway-ca.crt"
  health:
    listen_address: "0.0.0.0:8484"
`
