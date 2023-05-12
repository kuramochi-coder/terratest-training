package test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/kuramochi-coder/terratest-training/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAksIntegration(t *testing.T) {
	t.Parallel()

	// Declare the expected values.
	expectedRgName := "cri3-aks01_re1-ez-sit-004"
	expectedClusterName := "cri3-aks01-ez-cls_r1-sit-004"

	// NameSpace to use.
	myNameSpace := "aksinternetns-rk"

	// Get AKS credentials.
	// This will run `az aks get-credentials --resource-group RESOURCE_GROUP --name CLUSTER_NAME`
	// and fail the test if there are any errors.
	// Note that this will overwrite the default kubectl config file at ~/.kube/config.
	getAksCredentialsOutput, getAksCredentialsErr := exec.Command("az", "aks", "get-credentials", "--admin", "--resource-group", expectedRgName, "--name", expectedClusterName, "--overwrite-existing").Output()
	require.NoError(t, getAksCredentialsErr, "Error getting AKS credentials: %s", getAksCredentialsOutput)

	// Path to the Kubernetes resource config we will test.
	kubeResourcePath, err := filepath.Abs("./nginx-deployment.yml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := strings.ToLower(myNameSpace)

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file.
	options := k8s.NewKubectlOptions("", "", namespaceName)

	// Create the namespace for this test.
	k8s.CreateNamespace(t, options, namespaceName)

	// cleanup
	t.Cleanup(func() {
		// Delete the namespace at the end of the test.
		k8s.DeleteNamespace(t, options, namespaceName)

		// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
		k8s.KubectlDelete(t, options, kubeResourcePath)
	})

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// This will wait up to 10 seconds for the service to become available, to ensure that we can access it.
	k8s.WaitUntilServiceAvailable(t, options, "nginx-service", 10, 20*time.Second)

	// Now we verify that the service will successfully boot and start serving requests
	service := k8s.GetService(t, options, "nginx-service")

	// Calling Sleep method
	time.Sleep(30 * time.Second)

	// Get the pod name.
	podNameOutput, getPodErr := exec.Command("kubectl", "get", "pods", "-n", namespaceName, "--no-headers", "-o", "custom-columns=\":metadata.name\"").Output()
	podName := string(podNameOutput)
	podName = utils.TrimQuotes(podName)

	// Get the pod IP.
	podIPOut, getIPErr := exec.Command("kubectl", "get", "pod", podName, "-n", namespaceName, "--template", "'{{.status.podIP}}'").Output()
	podIP := string(podIPOut)
	podIP = utils.TrimQuotes(podIP)

	// Get the Ingress IP.
	ingressIPOut, getIngressIPErr := exec.Command("kubectl", "get", "ingress", "-n", namespaceName, "--no-headers", "-o", "custom-columns=\":status.loadBalancer.ingress[0].ip\"").Output()
	ingressIP := string(ingressIPOut)
	ingressIP = utils.TrimQuotes(ingressIP)

	// Handle the error above.
	if getPodErr != nil || getIPErr != nil || getIngressIPErr != nil {
		// Set retries.
		retries := 3

		for i := 0; i < retries; i++ {
			// Calling Sleep method
			time.Sleep(30 * time.Second)

			// Get the pod name.
			podNameOutput, getPodErr = exec.Command("kubectl", "get", "pods", "-n", namespaceName, "--no-headers", "-o=custom-columns=\":metadata.name\"").Output()
			podName := string(podNameOutput)
			podName = utils.TrimQuotes(podName)

			// Get the pod IP.
			podIPOut, getIPErr = exec.Command("kubectl", "get", "pod", podName, "-n", namespaceName, "--template", "'{{.status.podIP}}'").Output()
			podIP = string(podIPOut)
			podIP = utils.TrimQuotes(podIP)

			// Get the Ingress IP.
			ingressIPOut, getIngressIPErr = exec.Command("kubectl", "get", "ingress", "-n", namespaceName, "--no-headers", "-o=custom-columns=\":status.loadBalancer.ingress[0].ip\"").Output()
			ingressIP = string(ingressIPOut)
			ingressIP = utils.TrimQuotes(ingressIP)

			// Check if the error is nil.
			if getPodErr == nil && getIPErr == nil && getIngressIPErr == nil {
				break
			}
		}

		// Check if there are any errors
		if getPodErr != nil || getIPErr != nil || getIngressIPErr != nil {
			// Log the error and exit the program.
			t.Fatal(getPodErr, getIPErr, getIngressIPErr)
		}
	}

	// Test the kubernetes service is running.
	t.Run("TestAksService", func(t *testing.T) {
		t.Parallel()

		// Test that the service and endpoint is available.
		assert.NotNil(t, service)
		assert.NotNil(t, podIP)
	})

	// Test for AKS Connection within App tier.
	t.Run("TestAksConnectionWithinAppTier", func(t *testing.T) {
		t.Parallel()

		// Test the podIP for up to 5 minutes. This will only fail if we timeout waiting for the service to return a 200
		// response.
		podIPAddress := "http://" + podIP
		curlOutput, curlErr := exec.Command("curl", "-o", "-k", podIPAddress).Output()

		assert.NotNil(t, string(curlOutput))
		assert.Nil(t, curlErr)
	})

	// Test for AKS Connection from web tier.
	t.Run("TestAksConnectionFromWebTier", func(t *testing.T) {
		t.Parallel()

		// Test the podIP for up to 5 minutes. This will only fail if we timeout waiting for the service to return a 200
		// response.
		ingressIPAddress := "http://" + ingressIP
		curlOutput, curlErr := exec.Command("curl", "-o", "-k", ingressIPAddress).Output()

		assert.NotNil(t, string(curlOutput))
		assert.Nil(t, curlErr)
	})
}
