package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/unmarshall/explore-azgo-sdk/common"
	"k8s.io/utils/pointer"
)

func main() {
	connectConfig := common.GetAzureConnectConfig()
	common.DieOnError(connectConfig.Validate(), "invalid connect config")
	clients, err := common.NewClients(connectConfig)
	common.DieOnError(err, "failed to create clients")

	ctx := context.Background()
	client := clients.ResourceGraphClient
	//queryTemplate := `
	//Resources
	//| where type =~ 'Microsoft.Compute/virtualMachines'
	//| where resourceGroup =~ '%s'
	//| mv-expand bagexpansion=array tags
	//| where isnotempty(tags)
	//| where tags[0] startswith 'kubernetes.io-cluster-' or tags[0] startswith 'kubernetes.io-role-'
	//| distinct name
	//`
	queryTemplate := `
	Resources
	| where type =~ 'Microsoft.Compute/virtualMachines'
	| where resourceGroup =~ '%s'
	| where bag_keys(tags) hasprefix "kubernetes.io-cluster-"
	| where bag_keys(tags) hasprefix "kubernetes.io-role-"
	| project name
	`
	queryStr := fmt.Sprintf(queryTemplate, "shoot--mb-garden--sdktest")
	//fmt.Println(queryStr)
	//query := "Resources | where type =~ 'Microsoft.Compute/virtualMachines' | where resourceGroup =~ 'shoot--mb-garden--sdktest' | mv-expand bagexpansion=array tags | where isnotempty(tags) | where tags[0] startswith 'kubernetes.io-cluster-' or tags[0] startswith 'kubernetes.io-role-' | project name, type, tags"
	resp, err := client.Resources(ctx,
		armresourcegraph.QueryRequest{
			//Query: pointer.String("Resources | where type =~ 'Microsoft.Compute/virtualMachines' and tags[\"kubernetes.io-role-node\"] =~ '1' and tags[\"kubernetes.io-cluster-shoot--mb-garden--sdktest\"] =~ '1' | where resourceGroup =~ 'shoot--mb-garden--sdktest' | project name, tags"),
			Query: to.Ptr(queryStr),
			Options: &armresourcegraph.QueryRequestOptions{
				ResultFormat: to.Ptr(armresourcegraph.ResultFormatObjectArray),
				//ResultFormat: to.Ptr(armresourcegraph.ResultFormatTable),
			},
			Subscriptions: []*string{pointer.String(connectConfig.SubscriptionID)},
		}, nil)

	common.DieOnError(err, "failed to query")
	fmt.Printf(" Resources found: %d\n", *resp.TotalRecords)
	fmt.Printf("resource.Data: %+v\n", resp.Data)
	fmt.Printf("facets: %+v\n", resp.Facets)
	if m, ok := resp.Data.([]interface{}); ok {
		for _, r := range m {
			items := r.(map[string]interface{})
			for k, v := range items {
				if sv, ok := v.(string); ok {
					fmt.Printf("string-value: %s\n", sv)
				}
				fmt.Printf("k: %s, v: %+v\n", k, v)
			}
		}
	}
}
