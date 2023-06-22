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
	client := clients.DiskClient

	ctx := context.Background()
	//diskName := vmName + "-os-disk"

	resp, err := client.Get(ctx, resourceGroup, "mb-1", nil)
	common.DieOnError(err, "failed to get the disk")
	respJsonBytes, err := resp.MarshalJSON()
	common.DieOnError(err, "failed to marshal disk GET response")
	fmt.Printf("Disk: %s", string(respJsonBytes))

}
