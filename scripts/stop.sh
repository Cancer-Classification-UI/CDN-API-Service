#!/bin/bash
echo "Killing and removing container..."
docker kill cdn-api
docker rm cdn-api
echo "Container stopped"
