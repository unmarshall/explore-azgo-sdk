package main

import (
	"context"
	"fmt"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

func main() {
	connectConfig := common.AzureConnectConfigFromEnv()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	//filterOnTags := map[string]string{
	//	//"worker.gardener.cloud_pool": "worker-blu9f",
	//	"Name": "shoot--mcm-ci--az-oot-target-worker-1-z2-866f6-ltl8n",
	//}

	vms, err := clients.ListVMs(context.Background(), nil)
	common.DieOnError(err, "failed to list VMs")
	for _, vm := range vms {
		fmt.Printf("ID: %s, name: %s, type: %s, Location: %s\n", *vm.ID, *vm.Name, *vm.Type, *vm.Location)
		fmt.Println("**********************Tags*********************")
		for k, v := range vm.Tags {
			fmt.Printf("key: %s, value: %s\n", k, *v)
		}
		fmt.Println("-------------------------------------------------------------------------------------------")
	}
}
