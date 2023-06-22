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
	client := clients.DiskClient

	ctx := context.Background()
	diskName := "mb-1"
	poller, err := client.BeginDelete(ctx, resourceGroup, diskName, nil)
	common.DieOnError(err, "failed to delete disk")
	_, err = poller.PollUntilDone(ctx, nil)
	common.DieOnError(err, "poll failed waiting to delete the disk")
	fmt.Println("Successfully delete the disk")

}
