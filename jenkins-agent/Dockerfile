# Start with the official Docker image which already has the Docker CLI and Git
FROM docker:latest

# Install necessary tools like curl and bash
RUN apk add --no-cache curl bash

# Download the official Helm installer script and execute it
# This is the most reliable method recommended by the Helm project
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && \
    ./get_helm.sh

# Verify the installation to ensure it worked correctly
RUN helm version