#!/bin/bash

# Example script to test the Go API
# Make sure the server is running: go run main.go

API_URL="http://localhost:8080/api/v1"

echo "🚀 Testing Go API - Scalable Architecture"
echo "==========================================="
echo ""

# Health Check
echo "1️⃣  Health Check"
curl -s "$API_URL/health" | jq '.'
echo ""

# Create Items
echo "2️⃣  Creating Items"
ITEM1=$(curl -s -X POST "$API_URL/items" \
  -H "Content-Type: application/json" \
  -d '{"name":"MacBook Pro","description":"16-inch, M3 Max"}')
echo "Created Item 1:"
echo "$ITEM1" | jq '.'
ITEM1_ID=$(echo "$ITEM1" | jq -r '.id')
echo ""

ITEM2=$(curl -s -X POST "$API_URL/items" \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","description":"Pro Max, 256GB"}')
echo "Created Item 2:"
echo "$ITEM2" | jq '.'
ITEM2_ID=$(echo "$ITEM2" | jq -r '.id')
echo ""

# Create Clients
echo "3️⃣  Creating Clients"
CLIENT1=$(curl -s -X POST "$API_URL/clients" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith","email":"alice@example.com","phone":"555-0101"}')
echo "Created Client 1:"
echo "$CLIENT1" | jq '.'
CLIENT1_ID=$(echo "$CLIENT1" | jq -r '.id')
echo ""

CLIENT2=$(curl -s -X POST "$API_URL/clients" \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob Johnson","email":"bob@example.com","phone":"555-0202"}')
echo "Created Client 2:"
echo "$CLIENT2" | jq '.'
CLIENT2_ID=$(echo "$CLIENT2" | jq -r '.id')
echo ""

# List All Items
echo "4️⃣  Listing All Items"
curl -s "$API_URL/items" | jq '.'
echo ""

# List All Clients
echo "5️⃣  Listing All Clients"
curl -s "$API_URL/clients" | jq '.'
echo ""

# Get Single Item
echo "6️⃣  Getting Item by ID: $ITEM1_ID"
curl -s "$API_URL/items/$ITEM1_ID" | jq '.'
echo ""

# Get Single Client
echo "7️⃣  Getting Client by ID: $CLIENT1_ID"
curl -s "$API_URL/clients/$CLIENT1_ID" | jq '.'
echo ""

# Update Item
echo "8️⃣  Updating Item"
curl -s -X PUT "$API_URL/items/$ITEM1_ID" \
  -H "Content-Type: application/json" \
  -d '{"name":"MacBook Pro (Updated)","description":"16-inch, M3 Max, 64GB RAM"}' | jq '.'
echo ""

# Update Client
echo "9️⃣  Updating Client"
curl -s -X PUT "$API_URL/clients/$CLIENT1_ID" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith (Updated)","email":"alice.smith@example.com","phone":"555-9999"}' | jq '.'
echo ""

# Delete Item
echo "🔟 Deleting Item: $ITEM2_ID"
curl -s -X DELETE "$API_URL/items/$ITEM2_ID" -w "\nStatus: %{http_code}\n"
echo ""

# Verify Deletion
echo "1️⃣1️⃣  Verifying Item Deleted - Listing Remaining Items"
curl -s "$API_URL/items" | jq '.'
echo ""

echo "✅ Test Complete!"
echo ""
echo "💡 Tips:"
echo "  - Add more resources by following the pattern in models/, handlers/, router/"
echo "  - The storage layer uses generics - works with any model type"
echo "  - Switch to PostgreSQL by implementing storage.Store[T] interface"
