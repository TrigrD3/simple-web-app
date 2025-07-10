# Start with the official Docker image which already has the Docker CLI and Git
FROM docker:latest

# Set environment variable for Helm version
ENV HELM_VERSION=v3.9.4

# Install necessary tools like curl and bash, then download and install Helm
RUN apk add --no-cache curl bash && \
    curl -sSL https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar xz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf linux-amd64

# Verify installation
RUN helm version