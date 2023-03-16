# terratest-training
Terratest training for everyone

# Getting Started
## Install Go
```bash
brew update && brew install go
```
## Install Azure CLI
```bash
brew install azure-cli
```
## Login to Azure
```bash
az login
```
## List Subscriptions
```bash
az account list
```
## Set the Subscription ID
```bash
az account set -s <your_subscription_id>
```
## Environment Variables Setup
```bash
export ARM_SUBSCRIPTION_ID=<your_subscription_id>
```

# Running the Tests
## Run the Example Test
```bash
go test -v test/azure/terraform_azure_example_test.go
```
