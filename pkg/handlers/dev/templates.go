package dev

type templateParams struct {
	TenantID     string
	PolicyName   string
	PolicyID     string
	DiscoveryURL string
	TenantKey    string
	ControlPlane struct {
		Enabled        bool
		Address        string
		ClientCertPath string
		ClientKeyPath  string
	}
	DecisionLogging bool
	DecisionLogger  struct {
		EMSAddress     string
		StorePath      string
		ClientCertPath string
		ClientKeyPath  string
	}
	SeedMetadata bool
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
      aserto-discovery:
        url: {{ .DiscoveryURL }}
        credentials:
          bearer:
            token: "{{ .TenantKey }}"
            scheme: "basic"
        headers:
          Aserto-Tenant-Id: {{ .TenantID }}
    discovery:
      service: aserto-discovery
      resource: {{ .PolicyName }}/{{ .PolicyName }}/opa
{{ if .ControlPlane.Enabled }}
controller:
  enabled: true
  server:
    address: {{ .ControlPlane.Address }}
    client_cert_path: {{ .ControlPlane.ClientCertPath }}
    client_key_path: {{ .ControlPlane.ClientKeyPath }}
{{ else }}
controller:
  enabled: false
{{ end }}
{{ if .DecisionLogging }}
decision_logger:
  type: self
  config:
    store_directory: {{ .DecisionLogger.StorePath }}
    scribe:
      address: {{ .DecisionLogger.EMSAddress }}
      client_cert_path: {{ .DecisionLogger.ClientCertPath }}
      client_key_path: {{ .DecisionLogger.ClientKeyPath }}
      ack_wait_seconds: 30
      headers:
        Aserto-Tenant-Id: {{ .TenantID }}
    shipper:
      publish_timeout_seconds: 2
{{ end }}
`

const configTemplateLocal = templatePreamble + `
opa:
  instance_id: "-"
  graceful_shutdown_period_seconds: 2
  local_bundles:
    paths: []
    skip_verification: true
`

const templatePreamble = `---
version: 2

logging:
  prod: true
  log_level: info

directory:
  db_path: ${TOPAZ_DIR}/db/directory.db
  seed_metadata: {{ .SeedMetadata }}

# remote directory is used to resolve the identity for the authorizer.
remote_directory:
  address: "0.0.0.0:9292" # set as default, it should be the same as the reader as we resolve the identity from the local directory service.
  insecure: true

# default jwt validation configuration
# jwt:
#   acceptable_time_skew_seconds: 5

api:
  health:
    listen_address: "0.0.0.0:9494"
  services:
    reader:
      grpc:
        listen_address: "0.0.0.0:9292"
        # if certs are not specified default certs will be generate with the format reader_grpc.*
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/grpc.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/grpc.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/grpc-ca.crt"
      gateway:
        listen_address: "0.0.0.0:9393"
        # allowed_origins include localhost by default
        allowed_origins:
        - https://*.aserto.com
        - https://*aserto-console.netlify.app
        # if certs are not specified the gateway will have the http: true flag enabled
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/gateway.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/gateway.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/gateway-ca.crt"      
    writer:
      grpc:
        listen_address: "0.0.0.0:9292"
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/grpc.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/grpc.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/grpc-ca.crt"
      gateway:
        listen_address: "0.0.0.0:9393"
        allowed_origins:
        - https://*.aserto.com
        - https://*aserto-console.netlify.app
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/gateway.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/gateway.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/gateway-ca.crt"      
    exporter:
      grpc:
        listen_address: "0.0.0.0:9292"
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/grpc.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/grpc.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/grpc-ca.crt"
      gateway:
        listen_address: "0.0.0.0:9393"
        allowed_origins:
        - https://*.aserto.com
        - https://*aserto-console.netlify.app
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/gateway.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/gateway.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/gateway-ca.crt"
    importer:
      grpc:
        listen_address: "0.0.0.0:9292"
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/grpc.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/grpc.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/grpc-ca.crt"
      gateway:
        listen_address: "0.0.0.0:9393"
        allowed_origins:
        - https://*.aserto.com
        - https://*aserto-console.netlify.app
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/gateway.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/gateway.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/gateway-ca.crt"
    
    authorizer:
      needs:
        - reader
      grpc:
        connection_timeout_seconds: 2
        listen_address: "0.0.0.0:8282"
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/grpc.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/grpc.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/grpc-ca.crt"
      gateway:
        listen_address: "0.0.0.0:8383"
        allowed_origins:
        - https://*.aserto.com
        - https://*aserto-console.netlify.app
        certs:
          tls_key_path: "${TOPAZ_DIR}/certs/gateway.key"
          tls_cert_path: "${TOPAZ_DIR}/certs/gateway.crt"
          tls_ca_cert_path: "${TOPAZ_DIR}/certs/gateway-ca.crt"
`
