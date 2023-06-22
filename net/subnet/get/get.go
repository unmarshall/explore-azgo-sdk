package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup      = "shoot--mb-garden--sdktest"
	virtualNetworkName = "shoot--mb-garden--sdktest1"
	subnetName         = "shoot--mb-garden--sdktest-nodes"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	nwClient := clients.SubnetClient
	resp, err := nwClient.Get(context.Background(), resourceGroup, virtualNetworkName, subnetName, nil)
	var respErr *azcore.ResponseError
	if err != nil {
		if errors.As(err, &respErr) {
			fmt.Printf("StatusCode: %d, ErrorCode: %s, Header: %+v", respErr.StatusCode, respErr.ErrorCode, respErr.RawResponse.Header)
		}
	}
	common.DieOnError(err, "failed to get subnet")
	fmt.Printf("Response.Name: %s, Response.ID: %s\n", *resp.Name, *resp.ID)
	fmt.Printf("Name: %s, ID: %s, Type: %s\n", *resp.Subnet.Name, *resp.Subnet.ID, *resp.Subnet.Type)
	propertyBytes, err := resp.Subnet.Properties.MarshalJSON()
	common.DieOnError(err, "failed to marshal subnet properties")
	fmt.Printf("Properties: %s", string(propertyBytes))
}
