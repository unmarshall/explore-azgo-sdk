package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
)

func main() {

	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")
	client := clients.InterfacesClient

	ctx := context.Background()
	pager := client.NewListPager(resourceGroup, nil)

	var pageCounter int32
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get next page, current page: %d: %v", pageCounter, err)
			break
		}
		for _, nic := range page.Value {
			fmt.Printf("Name: %s, ID: %s, Type: %s, Location: %s\n", *nic.Name, *nic.ID, *nic.Type, *nic.Location)
		}
		pageCounter++
	}

}
