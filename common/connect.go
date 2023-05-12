package common

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ErrMissingTenantID       = errors.New("tenant ID is missing")
	ErrMissingClientID       = errors.New("client ID is missing")
	ErrMissingClientSecret   = errors.New("client secret is missing")
	ErrMissingSubscriptionID = errors.New("subscription ID is missing")
)

type AzureConnectConfig struct {
	TenantID       string `json:"tenantID"`
	ClientID       string `json:"clientID"`
	ClientSecret   string `json:"clientSecret"`
	SubscriptionID string `json:"subscriptionID"`
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

func NewAzureConnectConfig(kubeConfigPath string) AzureConnectConfig {
	client := createClient(kubeConfigPath)
	secretsClient := client.CoreV1().Secrets("garden-core")
	secret, err := secretsClient.Get(context.Background(), "shoot-operator-az-team", v1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			fmt.Printf("secret is not found, check the secret name and namespace")
		}
		DieOnError(err, "failed to get secret for extracting azure credentials")
	}
	return createAzureConnectConfigFromSecret(secret)
}

func DefaultAzureConnectConfig() AzureConnectConfig {
	return NewAzureConnectConfig("virtual-garden-kubeconfig.yaml")
}

func createAzureConnectConfigFromSecret(secret *corev1.Secret) AzureConnectConfig {
	clientID := string(secret.Data["clientID"])
	clientSecret := string(secret.Data["clientSecret"])
	tenantID := string(secret.Data["tenantID"])
	subscriptionID := string(secret.Data["subscriptionID"])

	return AzureConnectConfig{
		TenantID:       tenantID,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		SubscriptionID: subscriptionID,
	}
}

func createClient(kubeConfigPath string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	DieOnError(err, "failed to load kubeconfig")
	clientSet, err := kubernetes.NewForConfig(config)
	DieOnError(err, "failed to create ClientSet")
	return clientSet
}
