package main

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
)

func main() {
	tenantID := os.Getenv("TENANT_ID")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	tokenCredential, err := createTokenCredential(tenantID, clientID, clientSecret)
	if err != nil {
		log.Fatalf("failed to create token credentials: %v", err)
	}
	log.Printf("tokenCredential: %v", tokenCredential)
	computeClientFactory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		log.Fatalf("failed to create compute client factory: %v", err)
	}
	vmClient := computeClientFactory.NewVirtualMachinesClient()
	res, err := vmClient.Get(context.Background(), "<resource-group-name>", "<vm-name>", nil)
	if err != nil {
		log.Fatalf("error getting vm: %v", err)
		return
	}
	log.Printf("vm: %v", *res.VirtualMachine.Name)
	log.Printf("vm: %v", *res.VirtualMachine.ID)
}

func createTokenCredential(tenantID string, clientID string, clientSecret string) (azcore.TokenCredential, error) {
	return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
}
