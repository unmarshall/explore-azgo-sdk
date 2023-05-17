package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	fmt.Printf("created VM: %v\n", createdVM)

}

func createVMParameters() armcompute.VirtualMachine {
	plan := armcompute.Plan{
		Name:      pointer.String(planName),
		Product:   pointer.String(product),
		Publisher: pointer.String(publisher),
	}

	vmSize := armcompute.VirtualMachineSizeTypesStandardA4V2
	helloCustomScript := "IyEvdXNyL2Jpbi9lbnYgYmFzaAoKZWNobyAiaGVsbG8i"

	vmProperties := armcompute.VirtualMachineProperties{
		HardwareProfile: &armcompute.HardwareProfile{
			VMSize: &vmSize,
		},
		NetworkProfile: &armcompute.NetworkProfile{
			NetworkAPIVersion:              nil,
			NetworkInterfaceConfigurations: nil,
			NetworkInterfaces:              nil,
		},
		OSProfile: &armcompute.OSProfile{
			AdminUsername: pointer.String(adminUserName),
			ComputerName:  pointer.String(vmName),
			CustomData:    &helloCustomScript,
			LinuxConfiguration: &armcompute.LinuxConfiguration{
				DisablePasswordAuthentication: pointer.Bool(true),
				SSH: &armcompute.SSHConfiguration{
					PublicKeys: []*armcompute.SSHPublicKey{
						createSSHPublicKey(),
					},
				},
			},
			RequireGuestProvisionSignal: nil,
			Secrets:                     nil,
			WindowsConfiguration:        nil,
		},
		StorageProfile: &armcompute.StorageProfile{
			DataDisks:      createDataDisks(),
			ImageReference: createImageReference(),
			OSDisk: &armcompute.OSDisk{
				CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
				Caching:      to.Ptr(armcompute.CachingTypesNone),
				DeleteOption: to.Ptr(armcompute.DiskDeleteOptionTypesDelete),
				DiskSizeGB:   pointer.Int32(50),
				ManagedDisk: &armcompute.ManagedDiskParameters{
					StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardSSDLRS),
				},
				Name: to.Ptr(vmName + "-os-disk"),
			},
		},
		UserData:               nil,
		VirtualMachineScaleSet: nil,
		InstanceView:           nil,
		ProvisioningState:      nil,
		TimeCreated:            nil,
		VMID:                   nil,
	}

	return armcompute.VirtualMachine{
		Location:         pointer.String(location),
		ExtendedLocation: nil,
		Identity:         nil,
		Plan:             &plan,
		Properties:       &vmProperties,
		Tags:             nil,
		Zones:            nil,
		ID:               nil,
		Name:             pointer.String(vmName),
		Resources:        nil,
		Type:             nil,
	}
}

func createSSHPublicKey() *armcompute.SSHPublicKey {
	homeDir, err := os.UserHomeDir()
	authorizedKeysPath := filepath.Join(homeDir, ".ssh/authorized_keys")
	common.DieOnError(err, "failed to get home directory")
	sshPublicKeyPath := filepath.Join(homeDir, ".ssh/id_rsa.pub")
	var sshBytes []byte
	if _, err := os.Stat(sshPublicKeyPath); err != nil {
		common.DieOnError(err, "failed to open path: "+sshPublicKeyPath)
	}
	sshBytes, err = os.ReadFile(sshPublicKeyPath)
	common.DieOnError(err, "failed to read file at: "+sshPublicKeyPath)
	return &armcompute.SSHPublicKey{
		KeyData: to.Ptr(string(sshBytes)),
		Path:    to.Ptr(fmt.Sprintf(authorizedKeysPath)),
	}
}

func createDataDisks() []*armcompute.DataDisk {
	return []*armcompute.DataDisk{
		{
			CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesEmpty),
			Lun:          pointer.Int32(1),
			Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
			DeleteOption: to.Ptr(armcompute.DiskDeleteOptionTypesDelete),
			DiskSizeGB:   pointer.Int32(10),
			ManagedDisk: &armcompute.ManagedDiskParameters{
				StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardSSDLRS),
			},
			Name: to.Ptr("mb-1"),
		},
	}
}

func createImageReference() *armcompute.ImageReference {
	urn := "sap:gardenlinux:greatest:934.7.0"
	splits := strings.Split(urn, ":")
	publisher := splits[0]
	offer := splits[1]
	sku := splits[2]
	version := splits[3]
	return &armcompute.ImageReference{
		Publisher: &publisher,
		Offer:     &offer,
		SKU:       &sku,
		Version:   &version,
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
