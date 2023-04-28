package common

import (
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	ErrMissingTenantID       = errors.New("tenant ID is missing")
	ErrMissingClientID       = errors.New("client ID is missing")
	ErrMissingClientSecret   = errors.New("client secret is missing")
	ErrMissingSubscriptionID = errors.New("subscription ID is missing")
)

type AzureConnectConfig struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
}

func (a *AzureConnectConfig) CreateTokenCredential() (azcore.TokenCredential, error) {
	return azidentity.NewClientSecretCredential(a.TenantID, a.ClientID, a.ClientSecret, nil)
}

func (a *AzureConnectConfig) Validate() error {
	if len(strings.TrimSpace(a.TenantID)) == 0 {
		return ErrMissingTenantID
	}
	if len(strings.TrimSpace(a.ClientID)) == 0 {
		return ErrMissingClientID
	}
	if len(strings.TrimSpace(a.ClientSecret)) == 0 {
		return ErrMissingClientSecret
	}
	if len(strings.TrimSpace(a.SubscriptionID)) == 0 {
		return ErrMissingSubscriptionID
	}
	return nil
}
