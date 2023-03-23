package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAksResource(t *testing.T) {
	t.Parallel()
	SubscriptionID := ""
	// Declare the expected values.
	expectedRgName := "cri3-aks01_re1-ez-sit-004"
	expectedClusterName := "cri3-aks01-ez-cls_r1-sit-004"
	expectedNodeCount := 1

	// Test for resource presence.
	t.Run("TestResourcesExists", func(t *testing.T) {
		t.Parallel()
		// Test if resource group exists.
		exists := azure.ResourceGroupExists(t, expectedRgName, SubscriptionID)
		assert.True(t, exists, fmt.Sprintf("Resource group (%s) does not exist", expectedRgName))

		// Test if managed cluster exists.
		cluster, err := azure.GetManagedClusterE(t, expectedRgName, expectedClusterName, "")
		require.NoError(t, err)
		assert.NotNil(t, cluster, fmt.Sprintf("Managed cluster (%s) does not exist", *cluster.Name))
	})

	// Test for AKS Cluster.
	t.Run("TestAksCluster", func(t *testing.T) {
		t.Parallel()
		// Look up the cluster node count.
		cluster, err := azure.GetManagedClusterE(t, expectedRgName, expectedClusterName, "")
		require.NoError(t, err)
		actualCount := *(*cluster.ManagedClusterProperties.AgentPoolProfiles)[0].Count

		// Test that the Node count matches the Terraform specification.
		assert.Equal(t, int32(expectedNodeCount), actualCount)
	})
}
