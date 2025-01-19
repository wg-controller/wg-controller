#!/bin/bash

# Load .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

GHCR_USERNAME="wg-controller"  # GitHub organization name
IMAGE_NAME="wg-controller-ui"  # Name of the Docker image
REPO="ghcr.io/$GHCR_USERNAME/$IMAGE_NAME"  # GHCR repository name

# Get the current git tag
IMAGE_TAG=$(git describe --tags --abbrev=0)

# If no tag is found, exit with an error
if [ -z "$IMAGE_TAG" ]; then
    echo "Error: No git tag found. Exiting."
    exit 1
fi

echo "Using image tag: $IMAGE_TAG"

# Authenticate to GHCR
echo $GHCR_TOKEN | docker login ghcr.io -u $GHCR_USERNAME --password-stdin

# Build the Docker image
docker buildx build --platform linux/amd64 --platform linux/arm64 --build-arg GHCR_TOKEN=$GHCR_TOKEN --build-arg IMAGE_TAG=$IMAGE_TAG --build-arg IMAGE_NAME=$IMAGE_NAME -t $REPO:$IMAGE_TAG -t $REPO:latest . --push

echo "Publish Script Complete"
