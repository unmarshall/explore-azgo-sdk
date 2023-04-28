package main

import (
	"context"
	"log"

	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-blu9f-z1-8f464-jb82h"
)

func main() {
	connectConfig := common.AzureConnectConfigFromEnv()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	vm, err := clients.GetVM(context.Background(), resourceGroup, vmName)
	common.DieOnError(err, "failed to get machine")
	log.Printf("vm: %v", vm.Name)
	log.Printf("vm: %v", vm.ID)
}
