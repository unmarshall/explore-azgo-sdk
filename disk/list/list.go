package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

func main() {

	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")
	client := clients.DiskClient

	ctx := context.Background()

	pager := client.NewListPager(nil)
	var pageCounter int32
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get next page, current page: %d: %v", pageCounter, err)
			break
		}
		for _, disk := range page.Value {
			fmt.Printf("Name: %s, ID: %s, Location: %s, Type: %s\n", *disk.Name, *disk.ID, *disk.Location, *disk.Type)
		}
		pageCounter++
	}

}
