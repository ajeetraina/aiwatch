#!/bin/sh
set -e

# Check for Docker socket and fix permissions if needed
if [ -e /var/run/docker.sock ]; then
  echo "Docker socket found at /var/run/docker.sock"
  
  # Get current permissions
  DOCKER_SOCKET_GID=$(stat -c '%g' /var/run/docker.sock)
  echo "Docker socket group ID: $DOCKER_SOCKET_GID"
  
  # Get docker group ID inside container
  DOCKER_GROUP_ID=$(getent group docker | cut -d: -f3)
  echo "Docker group ID in container: $DOCKER_GROUP_ID"
  
  # Check if we need to adjust the docker group ID to match socket
  if [ -n "$DOCKER_SOCKET_GID" ] && [ "$DOCKER_SOCKET_GID" != "$DOCKER_GROUP_ID" ]; then
    echo "Adjusting docker group ID to match socket"
    
    # We can't modify the group ID in Alpine without root, so we use a workaround
    # by creating a new group with the correct ID and adding our user to it
    addgroup -g "$DOCKER_SOCKET_GID" -S dockerhost
    addgroup appuser dockerhost
    
    echo "Fixed permissions for Docker socket access"
  fi
fi

# For debugging purposes
echo "Running as user $(id)"
echo "Current groups: $(id -G)"
docker --version || echo "Docker command not found or not working"
docker info || echo "Docker info command failed"

# Execute the original command
exec "$@"