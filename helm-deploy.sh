#!/bin/bash

export NAMESPACE="simple-web-app"
export CHART_DIR="/Users/arya/Downloads/Repo/simple-web-app/simple-web-app"  # Update this with the path to your local chart
export SERVICE_NAME="go-web-app"

echo "deploying $SERVICE_NAME"
echo "$NAMESPACE"

# Remove the helm-s3 plugin installation and repo-related commands
# helm plugin install https://github.com/hypnoglow/helm-s3.git --version 0.13.0
# helm repo add $HELM_REPO_NAME $HELM_REPO_URL
# helm repo update
# helm repo list
# helm search repo $HELM_REPO_NAME

# Deploy using the local chart
helm upgrade --install ${SERVICE_NAME} "$CHART_DIR" --namespace ${NAMESPACE} --set "image.tag=v4" -f values.yaml

# helm uninstall ${SERVICE_NAME} --namespace ${NAMESPACE} 

