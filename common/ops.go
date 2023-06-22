package common

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/marketplaceordering/armmarketplaceordering"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Clients struct {
	ResourcesClient             *armresources.Client
	ResourceGroupClient         *armresources.ResourceGroupsClient
	VirtualNetworkClient        *armnetwork.VirtualNetworksClient
	SubnetClient                *armnetwork.SubnetsClient
	InterfacesClient            *armnetwork.InterfacesClient
	VirtualMachineClient        *armcompute.VirtualMachinesClient
	ImagesClient                *armcompute.ImagesClient
	VirtualMachineImagesClient  *armcompute.VirtualMachineImagesClient
	DiskClient                  *armcompute.DisksClient
	MarketPlaceAgreementsClient *armmarketplaceordering.MarketplaceAgreementsClient
	ResourceSKUClient           *armcompute.ResourceSKUsClient
	ResourceGraphClient         *armresourcegraph.Client
}

func NewClients(connectConfig AzureConnectConfig) (*Clients, error) {
	tokenCredential, err := connectConfig.CreateTokenCredential()
	if err != nil {
		return nil, err
	}
	clients := Clients{}
	clients.ResourcesClient, err = createResourcesClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create resources client")
	clients.ResourceGroupClient, err = createResourceGroupClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create resource group client")
	clients.VirtualMachineClient, err = createVirtualMachineClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create vm client")
	clients.ImagesClient, clients.VirtualMachineImagesClient, err = createImagesClients(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create image clients")
	clients.VirtualNetworkClient, clients.SubnetClient, clients.InterfacesClient, err = createNetworkClients(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create network clients")
	clients.MarketPlaceAgreementsClient, err = createMarketPlaceAgreementsClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create marketplace agreements client")
	clients.ResourceSKUClient, err = createResourceSKUsClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create resource SKU client")
	clients.ResourceGraphClient, err = createResourceGraphClient(tokenCredential)
	DieOnError(err, "failed to create resource graph client")
	clients.DiskClient, err = createDiskClient(connectConfig.SubscriptionID, tokenCredential)
	return &clients, nil
}

func (c *Clients) GetVM(ctx context.Context, resourceGroupName, vmName string) (armcompute.VirtualMachine, error) {
	res, err := c.VirtualMachineClient.Get(ctx, resourceGroupName, vmName, nil)
	if err != nil {
		return armcompute.VirtualMachine{}, err
	}
	return res.VirtualMachine, nil
}

func (c *Clients) ResourceGroupExists(ctx context.Context, resourceGroupName string) (bool, error) {
	resp, err := c.ResourceGroupClient.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (c *Clients) GetResourceGroup(ctx context.Context, resourceGroupName string) (*armresources.ResourceGroup, error) {
	resp, err := c.ResourceGroupClient.Get(ctx, resourceGroupName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.ResourceGroup, nil
}

func (c *Clients) ListVMs(ctx context.Context, tags map[string]string) ([]*armcompute.VirtualMachine, error) {
	var vms []*armcompute.VirtualMachine
	pager := c.VirtualMachineClient.NewListAllPager(&armcompute.VirtualMachinesClientListAllOptions{
		Filter: nil,
	})
	var pageCount int32
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Printf("failed to go to next page, current at page number: %d: %v", pageCount, err)
			return vms, err
		}
		for _, vm := range page.Value {
			vms = append(vms, vm)
		}
	}
	return vms, nil
}

func (c *Clients) ListResources(ctx context.Context, filters map[string]string) ([]*armresources.GenericResourceExpanded, error) {
	var resources []*armresources.GenericResourceExpanded
	filter := "resourceGroup eq 'shoot--mcm-ci--az-oot-target' and resourceType eq 'Microsoft.Compute/virtualMachines'"
	pager := c.ResourcesClient.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})
	var pageCount int32
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Printf("failed to go to the next page, currently at page number: %d: %v", pageCount, err)
			return resources, err
		}
		for _, r := range page.Value {
			resources = append(resources, r)
		}
	}
	return resources, nil
}

//func createContainsTagsFilter(tags map[string]string) *string {
//	if tags == nil || len(tags) == 0 {
//		return nil
//	}
//	var filter strings.Builder
//	var counter int
//	for k, v := range tags {
//		filter.WriteString(fmt.Sprintf("%s eq %s", k, v))
//		if counter < len(tags)-1 {
//			filter.WriteString(" and ")
//		}
//		counter++
//	}
//	str := filter.String()
//	return &str
//}

func createVirtualMachineClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armcompute.VirtualMachinesClient, error) {
	factory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewVirtualMachinesClient(), nil
}

func createImagesClients(subscriptionID string, tokenCredential azcore.TokenCredential) (*armcompute.ImagesClient, *armcompute.VirtualMachineImagesClient, error) {
	factory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, nil, err
	}
	return factory.NewImagesClient(), factory.NewVirtualMachineImagesClient(), nil
}

func createNetworkClients(subscriptionID string, tokenCredential azcore.TokenCredential) (*armnetwork.VirtualNetworksClient, *armnetwork.SubnetsClient, *armnetwork.InterfacesClient, error) {
	factory, err := armnetwork.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	return factory.NewVirtualNetworksClient(), factory.NewSubnetsClient(), factory.NewInterfacesClient(), nil
}

func createResourceGroupClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armresources.ResourceGroupsClient, error) {
	factory, err := armresources.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewResourceGroupsClient(), nil
}

func createResourcesClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armresources.Client, error) {
	factory, err := armresources.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewClient(), nil
}

func createMarketPlaceAgreementsClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armmarketplaceordering.MarketplaceAgreementsClient, error) {
	factory, err := armmarketplaceordering.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewMarketplaceAgreementsClient(), nil
}

func createResourceSKUsClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armcompute.ResourceSKUsClient, error) {
	factory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewResourceSKUsClient(), nil
}

func createResourceGraphClient(tokenCredential azcore.TokenCredential) (*armresourcegraph.Client, error) {
	factory, err := armresourcegraph.NewClientFactory(tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewClient(), nil
}

func createDiskClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armcompute.DisksClient, error) {
	factory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewDisksClient(), nil
}
