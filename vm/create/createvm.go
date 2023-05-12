package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/unmarshall/explore-azgo-sdk/common"
	"k8s.io/utils/pointer"
)

const (
	resourceGroup = "shoot--mb-garden--sdktest"
	vmName        = "shoot--mb-garden--sdktest-worker-bingo"
	location      = "westeurope"
	planName      = "greatest"
	product       = "gardenlinux"
	publisher     = "sap"
	adminUserName = "core"
)

var (
	clients *common.Clients
)

func main() {
	var err error
	connectConfig := common.AzureConnectConfigFromEnv()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err = common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	exists, err := vmExists()
	common.DieOnError(err, "failed to get vm")
	fmt.Printf("Vm Exists? %t", exists)

	parameters := createVMParameters()

	ctx := context.Background()

	vmClient := clients.VirtualMachineClient
	pollerResponse, err := vmClient.BeginCreateOrUpdate(ctx, resourceGroup, vmName, parameters, nil)
	common.DieOnError(err, "failed to make request to create vm")
	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	common.DieOnError(err, "polling for vm created failed")

	createdVM := resp.VirtualMachine

}

func createVMParameters() armcompute.VirtualMachine {
	plan := armcompute.Plan{
		Name:      pointer.String(planName),
		Product:   pointer.String(product),
		Publisher: pointer.String(publisher),
	}

	vmSize := armcompute.VirtualMachineSizeTypesStandardA4V2
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(secret.Data["userData"]))

	vmProperties := armcompute.VirtualMachineProperties{
		AdditionalCapabilities: nil,
		ApplicationProfile:     nil,
		AvailabilitySet:        nil,
		BillingProfile:         nil,
		CapacityReservation:    nil,
		DiagnosticsProfile:     nil,
		EvictionPolicy:         nil,
		ExtensionsTimeBudget:   nil,
		HardwareProfile: &armcompute.HardwareProfile{
			VMSize: &vmSize,
		},
		NetworkProfile: nil,
		OSProfile: &armcompute.OSProfile{
			AdminPassword:               nil,
			AdminUsername:               pointer.String(adminUserName),
			AllowExtensionOperations:    nil,
			ComputerName:                pointer.String(vmName),
			CustomData:                  nil,
			LinuxConfiguration:          nil,
			RequireGuestProvisionSignal: nil,
			Secrets:                     nil,
			WindowsConfiguration:        nil,
		},
		PlatformFaultDomain:     nil,
		Priority:                nil,
		ProximityPlacementGroup: nil,
		ScheduledEventsProfile:  nil,
		SecurityProfile:         nil,
		StorageProfile:          nil,
		UserData:                nil,
		VirtualMachineScaleSet:  nil,
		InstanceView:            nil,
		ProvisioningState:       nil,
		TimeCreated:             nil,
		VMID:                    nil,
	}

	return armcompute.VirtualMachine{
		Location:         pointer.String(location),
		ExtendedLocation: nil,
		Identity:         nil,
		Plan:             &plan,
		Properties:       nil,
		Tags:             nil,
		Zones:            nil,
		ID:               nil,
		Name:             pointer.String(vmName),
		Resources:        nil,
		Type:             nil,
	}
}

func vmExists() (bool, error) {
	_, err := clients.VirtualMachineClient.Get(context.Background(), resourceGroup, vmName, nil)
	var respErr *azcore.ResponseError
	if err != nil {
		if errors.As(err, &respErr) {
			if http.StatusNotFound == respErr.StatusCode {
				fmt.Printf("VM with name %s is not found", vmName)
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
