#!/bin/bash

# Function to run encoding and decoding tests
run_test() {
    secret="$1"
    total_shares="$2"
    threshold="$3"

    # Python encode
    python_encoded=$(python3 script.py split "$secret" "$total_shares" "$threshold")
    echo "Python Encoded: $python_encoded"

    # Go decode Python encoded
    go_decoded=$(go run script.go restore "$python_encoded")
    echo "Go Decoded: $go_decoded"

    # Check if the Go decoded matches the original secret
    if [ "$go_decoded" != "$secret" ]; then
        echo "Test failed: Go failed to decode Python encoded secret"
        exit 1
    fi

    # Go encode
    go_encoded=$(go run script.go split "$secret" "$total_shares" "$threshold")
    echo "Go Encoded: $go_encoded"

    # Python decode Go encoded
    python_decoded=$(python3 script.py restore "$go_encoded")
    echo "Python Decoded: $python_decoded"

    # Check if the Python decoded matches the original secret
    if [ "$python_decoded" != "$secret" ]; then
        echo "Test failed: Python failed to decode Go encoded secret"
        exit 1
    fi
}

# Test cases
declare -a secrets=(
    "Stefan"
    "I am a black cat and I have two heads and three legs"
    "My stinky feet do not stink as badly as yours, but it's only because you have no feet"
    "able place bird crane exhibit beach equip trap rigid squeeze judge swap segment bulk bargain never session frost buddy exit sunny coffee month short"
    "receive dune emerge need figure total poem convince dinner typical finger grit burger region search october ordinary vehicle wrong depart quick fall wheel police"
    "kangaroo chase tilt bike ecology favorite high glimpse satoshi country captain trade talk plate habit sun midnight library still script inside power crazy road"
)

# Thresholds
declare -a thresholds=("2 2" "2 3" "3 5" "2 7" "7 7")

# Iterate through test cases and thresholds
for secret in "${secrets[@]}"; do
    for threshold in "${thresholds[@]}"; do
        echo "Testing secret: '$secret' with threshold: $threshold"
        run_test "$secret" $threshold
        result=$(echo $?)
        if [ "$result" != "0" ]; then
            echo "Test failed. Stopping: $result"
            exit 1
        fi
        echo "Test passed for secret: '$secret' with threshold: $threshold"
        echo "----------------------------------------"
    done
done

echo "All tests passed successfully!"