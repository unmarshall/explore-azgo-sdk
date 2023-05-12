package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	offerType = armmarketplaceordering.OfferTypeVirtualmachine
	publisher = "sap"
	offer     = "gardenlinux"
	planName  = "greatest"
)

func main() {
	connectConfig := common.AzureConnectConfigFromEnv()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	client := clients.MarketPlaceAgreementsClient

	resp, err := client.Get(context.Background(), offerType, publisher, offer, planName, nil)
	common.DieOnError(err, "failed to get the agreement")
	fmt.Printf("Name: %s, ID: %s, Type: %s\n", *resp.Name, *resp.ID, *resp.Type)
	fmt.Printf("Plan: %s, Product: %s, Publisher: %s, Accepted: %t\n",
		*resp.AgreementTerms.Properties.Plan,
		*resp.AgreementTerms.Properties.Product,
		*resp.AgreementTerms.Properties.Publisher,
		*resp.AgreementTerms.Properties.Accepted,
	)
}
