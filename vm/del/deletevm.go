package main

import (
	"context"
	"fmt"
	"time"

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

	vmClient := clients.VirtualMachineClient
	ctx := context.Background()
	vmName := "shoot--mb-garden--sdktest-worker-bingo"
	start := time.Now()
	poller, err := vmClient.BeginDelete(ctx, resourceGroup, vmName, nil)
	common.DieOnError(err, "failed to delete VM")
	_, err = poller.PollUntilDone(ctx, nil)
	common.DieOnError(err, "poll for VM failed")
	fmt.Println("successfully deleted VM")
	fmt.Printf("Total time taken: %fs\n", time.Now().Sub(start).Seconds())

}
