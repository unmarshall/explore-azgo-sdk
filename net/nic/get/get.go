package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-bingo"
)

var (
	nicName = fmt.Sprintf("%s-nic-alpha", vmName)
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")
	client := clients.InterfacesClient
	resp, err := client.Get(context.Background(), resourceGroup, nicName, nil)
	common.DieOnError(err, "failed to get nic")
	fmt.Printf("Name: %s, ID: %s, Type: %s\n", *resp.Name, *resp.ID, *resp.Type)
	fmt.Printf("Interface.Name: %s, Interface.ID: %s, Interface.Type: %s\n", *resp.Interface.Name, *resp.Interface.ID, *resp.Interface.Type)
	fmt.Println("********  Tags *********")
	for k, v := range resp.Tags {
		fmt.Printf("Key : %s, Value: %s\n", k, *v)
	}
	fmt.Println("********  Interface Tags *********")
	for k, v := range resp.Interface.Tags {
		fmt.Printf("Key : %s, Value: %s\n", k, *v)
	}
	fmt.Println("********  Interface Properties *********")
	fmt.Printf("Properties: %v\n", resp.Interface.Properties)
}
