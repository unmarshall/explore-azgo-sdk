package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/unmarshall/explore-azgo-sdk/common"
	"k8s.io/utils/pointer"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	ctx := context.Background()
	client := clients.ResourceGraphClient
	resources, err := client.Resources(ctx,
		armresourcegraph.QueryRequest{
			Query:         pointer.String("Resources | where type =~ 'Microsoft.Compute/virtualMachines' | where resourceGroup =~ 'shoot--mb-garden--sdktest' | limit 3"),
			Subscriptions: []*string{pointer.String(connectConfig.SubscriptionID)},
		}, nil)
	common.DieOnError(err, "failed to query")
	fmt.Printf(" Resources found: %d", *resources.TotalRecords)
	fmt.Printf("resource.Data: %v", resources.Data)
}
