package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	// Assuming common.Env, common.DevDockerComposeEnv, common.AksStagingEnv, common.Info, common.Warn, common.InfoFields are accessible
	// If not, they might need to be passed or accessed differently, or logging simplified.
	// For zap logger (if common.InfoFields uses it):
	"go.uber.org/zap" // Ensure this import is present if using zap for InfoFields/Warn
	capi "github.com/hashicorp/consul/api"
)

type MySQLConfig struct {
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	User             string `json:"user,omitempty"`
	Pwd              string `json:"pwd,omitempty"` // Changed from "password" to "pwd" to match task
	Database         string `json:"database,omitempty"`
	DisableTLS       bool   `json:"disable_tls,omitempty"` // Used for local dev to bypass TLS
	ExternDataSource string `json:"extern_data_source,omitempty"`
}

// GetMySqlConfigFromConsul retrieves MySQL configuration from Consul.
// It now returns an error to allow for more granular error handling.
func GetMySqlConfigFromConsul(client *capi.Client, prefix string, key string) (*MySQLConfig, error) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		return nil, fmt.Errorf("Consul KV().Get failed for %s/%s: %w", prefix, key, err)
	}
	if val == nil || val.Value == nil {
		// Consider returning a more specific error or a nil config if key not found is acceptable in some cases
		return nil, fmt.Errorf("MySQL config key not found in Consul: %s/%s. Value is nil", prefix, key)
	}
	config := &MySQLConfig{}
	err = json.Unmarshal(val.Value, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal MySQL config from Consul for key %s/%s: %w", prefix, key, err)
	}
	return config, nil
}

// GetCurrentMySQLConfig determines if MySQL config should be loaded from ENV or Consul.
func GetCurrentMySQLConfig(client *capi.Client, consulPrefix string) (*MySQLConfig, error) {
	if Env.RunEnv == DevDockerComposeEnv {
		// Attempt to load from environment variables for Docker Compose
		host := os.Getenv("MYSQL_HOST")
		portStr := os.Getenv("MYSQL_PORT")
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD") // MYSQL_PASSWORD can be empty for local dev
		database := os.Getenv("MYSQL_DATABASE")
		disableTLSStr := os.Getenv("MYSQL_DISABLE_TLS")

		if host != "" && portStr != "" && user != "" && database != "" {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return nil, fmt.Errorf("invalid MYSQL_PORT ENV VAR: '%s': %w", portStr, err)
			}

			disableTLS := false
			if disableTLSStr == "true" {
				disableTLS = true
			}

			// Placeholder for ExternDataSource if needed for dev_docker_compose
			// externDataSourceEnv := os.Getenv("MYSQL_EXTERN_DATASOURCE")

			// Assuming InfoFields is available and works with zap.String etc.
			// If not, replace with standard logging, e.g., log.Printf or similar.
			InfoFields("Loaded MySQL config from ENV variables for dev_docker_compose",
				zap.String("host", host), zap.Int("port", port), zap.String("user", user), zap.String("database", database), zap.Bool("disableTLS", disableTLS))

			return &MySQLConfig{
				Host:       host,
				Port:       port,
				User:       user,
				Pwd:        password,
				Database:   database,
				DisableTLS: disableTLS,
				// ExternDataSource: externDataSourceEnv, // If read from ENV
			}, nil
		}
		// Log a warning if essential ENV VARS are missing in dev_docker_compose mode
		// Assuming Warn is available and works with zap.String etc.
		// If not, replace with standard logging.
		Warn("dev_docker_compose mode: Not all required MySQL ENV VARS (MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_DATABASE) are set. Attempting fallback to Consul.")
	}

	// Fallback to Consul for other environments or if ENV VARS were incomplete in dev_docker_compose
	var consulKey string
	// Assuming StagingEnv is defined in common.const or similar
	if Env.RunEnv == AksStagingEnv || Env.RunEnv == StagingEnv {
		consulKey = "mysqlStaging"
	} else {
		consulKey = "mysql" // Default key for prod and other dev environments
	}

	// Assuming InfoFields is available
	InfoFields("Loading MySQL config from Consul", zap.String("consulKey", consulKey), zap.String("consulPrefix", consulPrefix))
	cfg, err := GetMySqlConfigFromConsul(client, consulPrefix, consulKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get MySQL config from Consul (key: %s): %w", consulKey, err)
	}
	return cfg, nil
}
