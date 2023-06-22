package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-bingo"
	location      = "westeurope"
	vnetName      = "shoot--mb-garden--sdktest"
	subnetName    = "shoot--mb-garden--sdktest-nodes"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")
	client := clients.InterfacesClient
	ctx := context.Background()
	nicName := vmName + "-nic-alpha"

	tags := map[string]*string{
		"purpose": to.Ptr("testing-new-api"),
		"version": to.Ptr("test-alpha"),
	}

	parameters := armnetwork.Interface{
		Location: to.Ptr(location),
		Properties: &armnetwork.InterfacePropertiesFormat{
			EnableIPForwarding: to.Ptr(true),
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr(nicName),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						Subnet:                    getSubnet(ctx, clients.SubnetClient, subnetName),
					},
				},
			},
			EnableAcceleratedNetworking: to.Ptr(false),
		},
		Tags: tags,
		Name: to.Ptr(nicName),
	}
	start := time.Now()
	pollResp, err := client.BeginCreateOrUpdate(ctx, resourceGroup, nicName, parameters, nil)
	common.DieOnError(err, "request to create NIC failed")
	resp, err := pollResp.PollUntilDone(ctx, nil)
	common.DieOnError(err, "polling failed while creating NIC")
	fmt.Println("Successfully create NIC")
	fmt.Printf("Total time taken: %fs\n", time.Now().Sub(start).Seconds())
	interfaceJsonBytes, err := resp.Interface.MarshalJSON()
	common.DieOnError(err, "failed to marshal created NIC Interface")
	fmt.Printf("Created Interface: %s", string(interfaceJsonBytes))
}

func getSubnet(ctx context.Context, client *armnetwork.SubnetsClient, subnetName string) *armnetwork.Subnet {
	subnetResp, err := client.Get(ctx, resourceGroup, vnetName, subnetName, nil)
	common.DieOnError(err, "failed to get subnet")
	return &armnetwork.Subnet{ID: subnetResp.ID}
}
