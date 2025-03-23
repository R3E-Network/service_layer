#!/bin/bash
# Deploy Neo Oracle Service to Azure Container Instances with Azure Confidential Computing

set -e

# Configuration
RESOURCE_GROUP="neo-oracle-rg"
LOCATION="eastus"
CONTAINER_NAME="neo-oracle"
IMAGE_NAME="neoracleregistry.azurecr.io/neo-oracle:latest"
CPU="2"
MEMORY="4"
ENCLAVE_TYPE="aci-vbs"

# Create resource group if it doesn't exist
az group create --name $RESOURCE_GROUP --location $LOCATION

# Create Azure Container Registry if it doesn't exist
ACR_NAME="neoracleregistry"
if ! az acr show --name $ACR_NAME --resource-group $RESOURCE_GROUP &> /dev/null; then
  echo "Creating Azure Container Registry..."
  az acr create --name $ACR_NAME --resource-group $RESOURCE_GROUP --sku Standard --admin-enabled true
fi

# Build and push Docker image
echo "Building and pushing Docker image..."
docker build -t $IMAGE_NAME .
az acr login --name $ACR_NAME
docker push $IMAGE_NAME

# Create a managed identity for the container
IDENTITY_NAME="neo-oracle-identity"
if ! az identity show --name $IDENTITY_NAME --resource-group $RESOURCE_GROUP &> /dev/null; then
  echo "Creating managed identity..."
  az identity create --name $IDENTITY_NAME --resource-group $RESOURCE_GROUP
fi
IDENTITY_ID=$(az identity show --name $IDENTITY_NAME --resource-group $RESOURCE_GROUP --query id -o tsv)

# Deploy container with Confidential Computing
echo "Deploying container with Confidential Computing..."
az container create \
  --resource-group $RESOURCE_GROUP \
  --name $CONTAINER_NAME \
  --image $IMAGE_NAME \
  --cpu $CPU \
  --memory $MEMORY \
  --secure-environment \
  --enclave-type $ENCLAVE_TYPE \
  --registry-login-server "neoracleregistry.azurecr.io" \
  --registry-username $(az acr credential show --name $ACR_NAME --query "username" -o tsv) \
  --registry-password $(az acr credential show --name $ACR_NAME --query "passwords[0].value" -o tsv) \
  --ports 8080 \
  --assign-identity $IDENTITY_ID \
  --environment-variables \
    CONFIG_FILE=/app/config/config.json \
  --ip-address Public \
  --dns-name-label neo-oracle

# Get the public IP
PUBLIC_IP=$(az container show --resource-group $RESOURCE_GROUP --name $CONTAINER_NAME --query ipAddress.ip -o tsv)
FQDN=$(az container show --resource-group $RESOURCE_GROUP --name $CONTAINER_NAME --query ipAddress.fqdn -o tsv)

echo "Neo Oracle Service deployed successfully!"
echo "Public IP: $PUBLIC_IP"
echo "FQDN: $FQDN"
echo "API endpoint: http://$FQDN:8080" 