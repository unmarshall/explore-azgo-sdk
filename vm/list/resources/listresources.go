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

	resources, err := clients.ListResources(context.Background(), nil)
	if err != nil {
		return
	}
	for _, r := range resources {
		fmt.Printf("ID: %s, Name: %s, Type: %s\n", *r.ID, *r.Name, *r.Type)
		fmt.Println("*********************  Tags   *******************")
		for k, v := range r.Tags {
			fmt.Printf("%s:%s\n", k, *v)
		}
		fmt.Printf("-------------------------------------------------\n")
	}

}
