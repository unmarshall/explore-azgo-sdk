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

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")
	client := clients.InterfacesClient
	ctx := context.Background()
	nicName := vmName + "-nic-alpha"
	poller, err := client.BeginDelete(ctx, resourceGroup, nicName, nil)
	common.DieOnError(err, "failed to start delete of NIC")
	resp, err := poller.PollUntilDone(ctx, nil)
	common.DieOnError(err, "failed to delete NIC")
	fmt.Printf("Successfully deleted NIC :%v\n", resp)
}
