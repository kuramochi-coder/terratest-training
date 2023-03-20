package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestNetworkingResource(t *testing.T) {
	// Test the networking resources.
	t.Parallel()
	SubscriptionID := ""
	// Declare the expected values.
	expectedNsgRgName := "cri3-nw-bz-nw_rg-nprod-004"
	expectedNsgName := "cri3-nw-bz-bastion01_nsg-nprod-004"
	// Expected NSG rules count is final count of the number of custom rules applied by us + 6 default rules.
	// e.g. 10 + 6 (default), so 16 should be used.
	expectedNsgCount := 16
	expectedSubnetRange := "192.168.5.0/26"
	expectedVNetName := "cri3-bastion-nprd"
	expectedSubnetName := "AzureBastionSubnet"
	expectedVnetRgName := "htx-platform"

	// Test for resource presence.
	t.Run("TestResourcesExists", func(t *testing.T) {
		t.Parallel()
		// Check the Virtual Network exists.
		assert.True(t, azure.VirtualNetworkExists(t, expectedVNetName, expectedVnetRgName, SubscriptionID))

		// Check the Subnet exists.
		assert.True(t, azure.SubnetExists(t, expectedSubnetName, expectedVNetName, expectedVnetRgName, SubscriptionID))
	})

	// Integrated network resource tests.
	t.Run("TestVirtualNetworksSubnets", func(t *testing.T) {
		t.Parallel()
		// Get the Virtual Network Subnets, check the Subnet exists and has the expected Address Prefix.
		actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, expectedVNetName, expectedVnetRgName, SubscriptionID)
		assert.NotNil(t, actualVnetSubnets[expectedSubnetName])
		assert.Equal(t, expectedSubnetRange, actualVnetSubnets[expectedSubnetName])
	})

	// Test for NSG rules.
	t.Run("TestNsgRules", func(t *testing.T) {
		t.Parallel()
		// Check the NSG exists and has the expected number of rules.
		rules, err := azure.GetAllNSGRulesE(expectedNsgRgName, expectedNsgName, "")
		assert.NoError(t, err)
		assert.Equal(t, expectedNsgCount, len(rules.SummarizedRules))
	})
}
