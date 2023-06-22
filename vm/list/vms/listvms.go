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

	vms, err := clients.ListVMs(context.Background(), nil)
	common.DieOnError(err, "failed to list VMs")
	for _, vm := range vms {
		fmt.Printf("ID: %s, name: %s, type: %s, Location: %s\n", *vm.ID, *vm.Name, *vm.Type, *vm.Location)
		fmt.Println("**********************Tags*********************")
		if vm.Identity != nil {
			fmt.Printf("Identity: %+v", *vm.Identity)
		}
		//for k, v := range vm.Tags {
		//	fmt.Printf("key: %s, value: %s\n", k, *v)
		//}
		fmt.Println("-------------------------------------------------------------------------------------------")
	}
}
