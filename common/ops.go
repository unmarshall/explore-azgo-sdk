package common

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Clients struct {
	ResourceGroupClient  *armresources.ResourceGroupsClient
	VirtualNetworkClient *armnetwork.VirtualNetworksClient
	VirtualMachineClient *armcompute.VirtualMachinesClient
	DiskClient           *armcompute.DisksClient
}

func NewClients(connectConfig AzureConnectConfig) (*Clients, error) {
	tokenCredential, err := connectConfig.CreateTokenCredential()
	if err != nil {
		return nil, err
	}
	clients := Clients{}
	clients.ResourceGroupClient, err = createResourceGroupClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create resource group client")
	clients.VirtualMachineClient, err = createVirtualMachineClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create vm client")
	clients.VirtualNetworkClient, err = createVirtualNetworkClient(connectConfig.SubscriptionID, tokenCredential)
	DieOnError(err, "failed to create virtual network client")
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
	//filter := createContainsTagsFilter(tags)
	//filter := "Name eq 'shoot--dev--azure-test-3'" // does not work
	//filter := "tagName eq 'Name'" // does not work
	//filter := "tagname eq 'Name'" // does not work
	//filter := "$filter=tagName eq 'Name'" // does not work
	//v := url.Values{}
	//v.Set("$filter", "tagname eq 'Name'")
	//filter := v.Encode()
	//filter := "tagName eq 'Name' and tagValue eq 'shoot--dev--azure-test-3'"
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

func createContainsTagsFilter(tags map[string]string) *string {
	if tags == nil || len(tags) == 0 {
		return nil
	}
	var filter strings.Builder
	var counter int
	for k, v := range tags {
		filter.WriteString(fmt.Sprintf("%s eq %s", k, v))
		if counter < len(tags)-1 {
			filter.WriteString(" and ")
		}
		counter++
	}
	str := filter.String()
	return &str
}

type customTransporter struct {
}

func (c customTransporter) Do(req *http.Request) (*http.Response, error) {
	log.Printf("req = %+v", req)
	return nil, fmt.Errorf("transporter does not respond")
}

func createVirtualMachineClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armcompute.VirtualMachinesClient, error) {
	factory, err := armcompute.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewVirtualMachinesClient(), nil
}

func createVirtualNetworkClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armnetwork.VirtualNetworksClient, error) {
	factory, err := armnetwork.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewVirtualNetworksClient(), nil
}

func createResourceGroupClient(subscriptionID string, tokenCredential azcore.TokenCredential) (*armresources.ResourceGroupsClient, error) {
	factory, err := armresources.NewClientFactory(subscriptionID, tokenCredential, nil)
	if err != nil {
		return nil, err
	}
	return factory.NewResourceGroupsClient(), nil
}
