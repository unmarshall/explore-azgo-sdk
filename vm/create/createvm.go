package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
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
	zone          = "1"
)

func main() {
	var err error
	ctx := context.Background()

	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	vmClient := clients.VirtualMachineClient
	exists, err := vmExists(ctx, vmClient)
	common.DieOnError(err, "failed to get vm")
	fmt.Printf("Vm Exists? %t", exists)
	parameters := createVMParameters()
	start := time.Now()
	pollerResponse, err := vmClient.BeginCreateOrUpdate(ctx, resourceGroup, vmName, parameters, nil)
	common.DieOnError(err, "failed to make request to create vm")
	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	common.DieOnError(err, "polling for vm created failed")

	createdVM := resp.VirtualMachine
	vmJsonBytes, err := createdVM.MarshalJSON()
	common.DieOnError(err, "failed to marshal vm json")
	fmt.Printf("created VM: %s\n", string(vmJsonBytes))
	fmt.Printf("Total time taken: %fs\n", time.Now().Sub(start).Seconds())
}

func createVMParameters() armcompute.VirtualMachine {
	plan := armcompute.Plan{
		Name:      pointer.String(planName),
		Product:   pointer.String(product),
		Publisher: pointer.String(publisher),
	}

	vmSize := armcompute.VirtualMachineSizeTypesStandardA4V2
	helloCustomScript := "IyEvdXNyL2Jpbi9lbnYgYmFzaAoKZWNobyAiaGVsbG8i"

	nicID := "/subscriptions/82b44c79-a5d4-4d74-8ff8-8639e79c1c39/resourceGroups/shoot--mb-garden--sdktest/providers/Microsoft.Network/networkInterfaces/shoot--mb-garden--sdktest-worker-bingo-nic-alpha"

	vmProperties := armcompute.VirtualMachineProperties{
		HardwareProfile: &armcompute.HardwareProfile{
			VMSize: &vmSize,
		},
		NetworkProfile: &armcompute.NetworkProfile{
			NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
				{
					ID: to.Ptr(nicID),
					Properties: &armcompute.NetworkInterfaceReferenceProperties{
						//DeleteOption: to.Ptr(armcompute.DeleteOptionsDelete),
						DeleteOption: to.Ptr(armcompute.DeleteOptionsDetach),
						Primary:      to.Ptr(true),
					},
				},
			},
		},
		OSProfile: &armcompute.OSProfile{
			AdminUsername: pointer.String(adminUserName),
			ComputerName:  pointer.String(vmName),
			CustomData:    &helloCustomScript,
			LinuxConfiguration: &armcompute.LinuxConfiguration{
				DisablePasswordAuthentication: pointer.Bool(true),
			},
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
	}

	tags := map[string]*string{
		"Name": to.Ptr(vmName),
		//"kubernetes.io-role-sdk-test": to.Ptr("1"),
		"kubernetes.io_arch": to.Ptr("amd64"),
		//"node.kubernetes.io_role":     to.Ptr("node"),
		"kubernetes.io-cluster-shoot--mb-garden--sdktest": to.Ptr("1"),
	}

	return armcompute.VirtualMachine{
		Location:   pointer.String(location),
		Plan:       &plan,
		Properties: &vmProperties,
		Tags:       tags,
		Zones:      []*string{to.Ptr(zone)},
		Name:       pointer.String(vmName),
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

func vmExists(ctx context.Context, client *armcompute.VirtualMachinesClient) (bool, error) {
	_, err := client.Get(ctx, resourceGroup, vmName, nil)
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
