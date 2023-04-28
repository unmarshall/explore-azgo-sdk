package common

import (
	"log"
	"os"
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
