package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/unmarshall/explore-azgo-sdk/common"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-bingo"
)

func main() {
	var err error
	ctx := context.Background()

	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	vmClient := clients.VirtualMachineClient
	vmUpdateParams := createVMUpdateParams()
	poller, err := vmClient.BeginUpdate(ctx, resourceGroup, vmName, vmUpdateParams, nil)
	common.DieOnError(err, "failed to start VM update")

	res, err := poller.PollUntilDone(ctx, nil)
	common.DieOnError(err, "failed to finish update of VM")
	vmJsonBytes, err := res.MarshalJSON()
	common.DieOnError(err, "failed to unmarshall response JSON")
	fmt.Printf("Successfully update the VM, updated VM: %+s\n", string(vmJsonBytes))
}

func createVMUpdateParams() armcompute.VirtualMachineUpdate {
	nicID := "/subscriptions/82b44c79-a5d4-4d74-8ff8-8639e79c1c39/resourceGroups/shoot--mb-garden--sdktest/providers/Microsoft.Network/networkInterfaces/shoot--mb-garden--sdktest-worker-bingo-nic-alpha"
	return armcompute.VirtualMachineUpdate{
		Properties: &armcompute.VirtualMachineProperties{
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: to.Ptr(nicID),
						Properties: &armcompute.NetworkInterfaceReferenceProperties{
							DeleteOption: to.Ptr(armcompute.DeleteOptionsDelete),
						},
					},
				},
			},
			//StorageProfile: &armcompute.StorageProfile{
			//	DataDisks:          nil,
			//	DiskControllerType: nil,
			//	ImageReference:     nil,
			//	OSDisk:             nil,
			//},
		},
		//Tags: nil,
	}
}
