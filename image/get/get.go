package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	location  = "westeurope"
	publisher = "sap"
	offer     = "gardenlinux"
	skus      = "greatest"
	version   = "934.7.0"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	client := clients.VirtualMachineImagesClient
	resp, err := client.Get(context.Background(), location, publisher, offer, skus, version, nil)
	common.DieOnError(err, "failed to get image")

	fmt.Printf("Name: %s, ID: %s, Location: %s\n", *resp.Name, *resp.ID, *resp.Location)
	fmt.Printf("VirtualMachineImage.Name: %s, VirtualMachineImage.ID: %s, VirtualMachineImage.Location: %s\n", *resp.VirtualMachineImage.Name, *resp.VirtualMachineImage.ID, *resp.VirtualMachineImage.Location)
	fmt.Println("******** Tags ********")
	for k, v := range resp.VirtualMachineImage.Tags {
		fmt.Printf("Key: %s, Value: %s\n", k, *v)
	}
	purchasePlan := resp.VirtualMachineImage.Properties.Plan
	fmt.Printf("PurchasePlan: [Name: %s, Product: %s, Publisher: %s]\n", *purchasePlan.Name, *purchasePlan.Product, *purchasePlan.Publisher)
}
