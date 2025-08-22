#!/bin/bash

# Login Performance Testing Script
echo "🚀 Testing Login Performance Optimizations"
echo "========================================="

# Test configuration
SERVER_URL="http://localhost:9000"
LOGIN_ENDPOINT="/api/v1/users/login"
TEST_EMAIL="john.doe@example.com"
TEST_PASSWORD="Test@123"
NUM_TESTS=10

echo "📊 Running $NUM_TESTS login requests..."
echo "Server: $SERVER_URL$LOGIN_ENDPOINT"
echo "User: $TEST_EMAIL"
echo ""

# Create JSON payload
JSON_PAYLOAD=$(cat <<EOF
{
    "email": "$TEST_EMAIL",
    "password": "$TEST_PASSWORD"
}
EOF
)

# Arrays to store results
declare -a response_times
declare -a status_codes
total_time=0
success_count=0

echo "⏱️  Test Results:"
echo "Request | Status | Duration (ms) | Server Duration (ms)"
echo "--------|--------|---------------|--------------------"

for i in $(seq 1 $NUM_TESTS); do
    # Make the request and capture timing
    response=$(curl -s -w "%{http_code},%{time_total}" \
                   -X POST \
                   -H "Content-Type: application/json" \
                   -d "$JSON_PAYLOAD" \
                   "$SERVER_URL$LOGIN_ENDPOINT")
    
    # Extract HTTP status and timing
    status_code=$(echo "$response" | tail -c 10 | cut -d',' -f1)
    curl_time=$(echo "$response" | tail -c 10 | cut -d',' -f2)
    
    # Convert to milliseconds
    curl_time_ms=$(echo "scale=0; $curl_time * 1000" | bc)
    
    # Extract server-side duration from response
    server_duration=$(echo "$response" | head -c -10 | jq -r '.data.performance.duration_ms // "N/A"' 2>/dev/null)
    if [ "$server_duration" = "null" ] || [ "$server_duration" = "" ]; then
        server_duration="N/A"
    fi
    
    # Store results
    response_times+=($curl_time_ms)
    status_codes+=($status_code)
    total_time=$(echo "$total_time + $curl_time_ms" | bc)
    
    if [ "$status_code" -eq 200 ]; then
        success_count=$((success_count + 1))
        status_icon="✅"
    else
        status_icon="❌"
    fi
    
    printf "%7d | %6s | %13s | %18s %s\n" $i $status_code "${curl_time_ms}" "$server_duration" "$status_icon"
done

echo ""
echo "📈 Performance Summary:"
echo "======================"

# Calculate statistics
avg_time=$(echo "scale=2; $total_time / $NUM_TESTS" | bc)
success_rate=$(echo "scale=1; $success_count * 100 / $NUM_TESTS" | bc)

# Find min and max
min_time=${response_times[0]}
max_time=${response_times[0]}
for time in "${response_times[@]}"; do
    if [ $time -lt $min_time ]; then
        min_time=$time
    fi
    if [ $time -gt $max_time ]; then
        max_time=$time
    fi
done

echo "Total Requests: $NUM_TESTS"
echo "Successful: $success_count ($success_rate%)"
echo "Failed: $((NUM_TESTS - success_count))"
echo ""
echo "⏱️  Response Time Statistics:"
echo "  Average: ${avg_time}ms"
echo "  Fastest: ${min_time}ms"  
echo "  Slowest: ${max_time}ms"
echo ""

# Performance evaluation
if [ $(echo "$avg_time < 100" | bc) -eq 1 ]; then
    echo "🎉 EXCELLENT: Average response time under 100ms"
elif [ $(echo "$avg_time < 200" | bc) -eq 1 ]; then
    echo "✅ GOOD: Average response time under 200ms"  
elif [ $(echo "$avg_time < 500" | bc) -eq 1 ]; then
    echo "⚠️  ACCEPTABLE: Average response time under 500ms"
else
    echo "❌ NEEDS IMPROVEMENT: Average response time over 500ms"
fi

echo ""
echo "🎯 Optimization Recommendations:"
if [ $(echo "$avg_time > 200" | bc) -eq 1 ]; then
    echo "• Consider connection pooling optimization"
    echo "• Review database query efficiency"  
    echo "• Check network latency"
fi
if [ $(echo "$success_rate < 100" | bc) -eq 1 ]; then
    echo "• Investigate failed requests"
    echo "• Check server error logs"
fi

echo ""
echo "📋 Test completed at $(date)"