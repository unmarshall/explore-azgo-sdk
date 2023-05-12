package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup      = "shoot--mb-garden--sdktest"
	virtualNetworkName = "shoot--mb-garden--sdktest"
	subnetName         = "shoot--mb-garden--sdktest-nodes"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	nwClient := clients.SubnetClient
	resp, err := nwClient.Get(context.Background(), resourceGroup, virtualNetworkName, subnetName, nil)
	common.DieOnError(err, "failed to get subnet")
	fmt.Printf("Response.Name: %s, Response.ID: %s\n", *resp.Name, *resp.ID)
	fmt.Printf("Name: %s, ID: %s, Type: %s, Properties: %v\n", *resp.Subnet.Name, *resp.Subnet.ID, *resp.Subnet.Type, *resp.Subnet.Properties)
}
