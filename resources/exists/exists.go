package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

	cli := clients.ResourcesClient
	ctx := context.Background()
	var rawResp *http.Response
	ctxWithResp := runtime.WithCaptureResponse(ctx, &rawResp)
	existence, err := cli.CheckExistence(
		ctxWithResp,
		resourceGroup,
		"Microsoft.Network",
		"/",
		"networkInterfaces",
		"shoot--mb-garden--sdktest-worker-bingo-nic-alpha",
		"2023-03-01-preview",
		//"2021-04-01",
		//"2021-02-01",
		nil,
	)
	fmt.Printf("Raw response: %v", *rawResp)
	common.DieOnError(err, "failed to check existence")
	fmt.Printf("exists? %t\n", existence.Success)

}
