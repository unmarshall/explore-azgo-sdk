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

	client := clients.ResourceSKUClient
	pager := client.NewListPager(nil)
	ctx := context.Background()

	var pageCounter int32
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get next page, current page: %d: %v", pageCounter, err)
			break
		}
		for _, resourceSKU := range page.Value {
			if *resourceSKU.Name == "Standard_A4_v2" {
				fmt.Printf("Name: %s, Size: %s, Family: %s\n", safeDeref(resourceSKU.Name), safeDeref(resourceSKU.Size), safeDeref(resourceSKU.Family))
			}
		}
		pageCounter++
	}
}

func safeDeref(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
