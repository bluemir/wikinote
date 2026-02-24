#!/bin/sh
# Find available container runtime (docker or podman)
if command -v docker >/dev/null 2>&1; then
    echo docker
elif command -v podman >/dev/null 2>&1; then
    echo podman
else
    echo "ERROR: docker or podman not found" >&2
    exit 1
fi
