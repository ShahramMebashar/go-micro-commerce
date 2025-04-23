#!/bin/bash

# Number of requests to make
NUM_REQUESTS=${1:-10}

echo "Making $NUM_REQUESTS requests to the product service API..."
echo "This will show which container handles each request..."

for i in $(seq 1 $NUM_REQUESTS); do
  echo "Request $i:"
  curl -s http://localhost:8888/api/products/health | jq
  echo ""
  sleep 0.5
done
