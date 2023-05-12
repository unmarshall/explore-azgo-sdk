package common

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/json"
)

func AzureConnectConfigFromEnv() AzureConnectConfig {
	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	return AzureConnectConfig{
		TenantID:       tenantID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		SubscriptionID: subscriptionID,
	}
}

func DieOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func GetAzureConnectConfig() AzureConnectConfig {
	configPath := filepath.Join("/tmp", "azure-connect-config.json")
	if _, err := os.Stat(configPath); err == nil {
		return unmarshalToAzureConnectConfig(configPath)
	}
	config := DefaultAzureConnectConfig()
	jsonBytes, err := json.Marshal(config)
	DieOnError(err, "failed to marshal azure connect config")
	DieOnError(os.WriteFile(configPath, jsonBytes, 644), "failed to write azure connect config to path: "+configPath)
	return config
}

func unmarshalToAzureConnectConfig(path string) AzureConnectConfig {
	var azureConnectConfig AzureConnectConfig
	jsonBytes, err := os.ReadFile(path)
	DieOnError(err, "failed to read config file: "+path)
	err = json.Unmarshal(jsonBytes, &azureConnectConfig)
	DieOnError(err, "failed to unmarshal azure connect config from path: "+path)
	return azureConnectConfig
}
