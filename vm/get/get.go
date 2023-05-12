package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-blu9f-z1-8f464-bdjrm"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	res, err := clients.VirtualMachineClient.Get(context.Background(), resourceGroup, vmName, nil)
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		if http.StatusNotFound == respErr.StatusCode {
			fmt.Printf("VM with name %s is not found", vmName)
			return
		}
	}
	common.DieOnError(err, "failed to get machine")
	fmt.Printf("Name: %s\n", *res.Name)
	fmt.Printf("ID: %s\n", *res.ID)
	fmt.Printf("VirtualMachine.Name: %s\n", *res.VirtualMachine.Name)
	fmt.Printf("VM: %v\n", res.VirtualMachine)

}
