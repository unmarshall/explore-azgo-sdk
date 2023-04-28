package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if http.StatusNotFound == respErr.StatusCode {
			fmt.Printf("VM with name %s is not found", vmName)
			return
		}
	}
	common.DieOnError(err, "failed to get machine")
	log.Printf("vm: %v", vm.Name)
	log.Printf("vm: %v", vm.ID)
}
