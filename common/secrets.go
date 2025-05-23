package common

import (
	"encoding/json"
	"fmt"
	// "log" // If common.Fatal uses standard log as fallback
	capi "github.com/hashicorp/consul/api"
	// Assuming common.InfoFields, common.Warn, common.Fatal are accessible
	// For zap logger (if common.InfoFields uses it):
	"go.uber.org/zap"
)

type SecretsConfig struct {
	JWTSecret     string `json:"jwt_secret"` // Changed to match task
	Secret        string `json:"secret"`       // Changed to match task
	SecretStaging string `json:"secret_staging"`
	OrderToken    string `json:"orderToken,omitempty"` // Kept existing field
	// Add other secrets if defined in the Consul "secrets" key
}

// Global variable to hold secrets
var Secrets *SecretsConfig // Initialized as nil, will be set by InitSecretsFromConsul

// InitSecretsFromConsul populates the global Secrets variable from Consul.
func InitSecretsFromConsul(client *capi.Client, prefix string, key string) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		// Using fmt.Printf for critical early errors if logger isn't ready or to avoid circular deps
		fmt.Printf("ERROR: Consul KV().Get failed for secrets key %s/%s: %v\n", prefix, key, err)
		// Depending on strictness, might os.Exit(1) or common.Fatal()
		// For now, initialize Secrets to avoid nil pointer, but it will be empty/zeroed.
		Secrets = &SecretsConfig{}
		return
	}
	if val == nil || val.Value == nil {
		fmt.Printf("WARN: Secrets key not found in Consul or value is nil: %s/%s. Secrets will be empty/zeroed.\n", prefix, key)
		Secrets = &SecretsConfig{}
		return
	}

	localSecrets := &SecretsConfig{}
	err = json.Unmarshal(val.Value, localSecrets)
	if err != nil {
		fmt.Printf("ERROR: Failed to unmarshal secrets from Consul for key %s/%s: %v. Secrets will be empty/zeroed.\n", prefix, key, err)
		Secrets = &SecretsConfig{} // Ensure Secrets is not nil
		return
	}
	Secrets = localSecrets
	// Assuming InfoFields and zap are available and configured
	InfoFields("Successfully loaded secrets from Consul.", zap.String("key", prefix+"/"+key))

	// Original logic for setting JWTSecret based on RunEnv seems to be superseded by direct "jwt_secret"
	// from Consul or ENV var override in main.go.
	// If "jwt_secret" is NOT expected from Consul directly, this part might need adjustment
	// or be removed if JWT_SECRET env var is the sole source for dev_docker_compose
	// and Consul's "jwt_secret" field is used for other envs.
	// For now, I'll keep the original logic commented out to reflect the task's new struct.
	// if Env.RunEnv == AksProductionEnv {
	//  Secrets.JWTSecret = Secrets.Secret
	// } else {
	//  Secrets.JWTSecret = Secrets.SecretStaging
	// }
}
