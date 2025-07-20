#!/bin/bash

# Test script for Workout Tracker API

API_BASE="http://localhost:8080/api"

echo "Testing Workout Tracker API..."
echo "================================"

# Test 1: Get all workouts (should be empty initially)
echo "1. Testing GET /api/workouts"
curl -s "$API_BASE/workouts" | jq '.' 2>/dev/null || curl -s "$API_BASE/workouts"
echo -e "\n"

# Test 2: Create a new workout
echo "2. Testing POST /api/workouts"
WORKOUT_JSON='{
    "name": "Test Workout",
    "date": "2024-01-15",
    "duration": 60,
    "notes": "This is a test workout"
}'

RESPONSE=$(curl -s -X POST "$API_BASE/workouts" \
    -H "Content-Type: application/json" \
    -d "$WORKOUT_JSON")

echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo -e "\n"

# Extract workout ID from response (if jq is available)
if command -v jq &> /dev/null; then
    WORKOUT_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
    
    if [ "$WORKOUT_ID" != "null" ] && [ "$WORKOUT_ID" != "" ]; then
        echo "3. Testing GET /api/workouts/$WORKOUT_ID"
        curl -s "$API_BASE/workouts/$WORKOUT_ID" | jq '.' 2>/dev/null || curl -s "$API_BASE/workouts/$WORKOUT_ID"
        echo -e "\n"
    fi
fi

echo "4. Testing GET /api/workouts (should now have one workout)"
curl -s "$API_BASE/workouts" | jq '.' 2>/dev/null || curl -s "$API_BASE/workouts"
echo -e "\n"

echo "API testing complete!"
echo "====================="
echo "To test the web interface:"
echo "1. Start the server: go run cmd/server/main.go"
echo "2. Open your browser to: http://localhost:8080"
