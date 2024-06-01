#!/bin/bash

# Function to run encoding and decoding tests
run_test() {
    secret="$1"
    threshold="$2"
    total_shares="$3"

    # Python encode
    python_encoded=$(docker run --rm shamir-python:latest split "$secret" $threshold $total_shares)
    echo "Python Encoded: $python_encoded"

    # Go decode Python encoded
    go_decoded=$(docker run --rm shamir-go:latest-arm64 restore $python_encoded)
    go_decoded=$(echo $go_decoded | sed 's/^"//;s/"$//')
    echo "Go Decoded: $go_decoded"

    # Check if the Go decoded matches the original secret
    if [ "$go_decoded" != "$secret" ]; then
        echo "expected: $secret. got: $go_decoded"
        echo "Test failed: Go failed to decode Python encoded secret"
        exit 1
    fi

    # Go encode
    go_encoded=$(docker run --rm shamir-go:latest-arm64 split "$secret" $threshold $total_shares)
    echo "Go Encoded: $go_encoded"

    # Python decode Go encoded
    python_decoded=$(docker run --rm shamir-python:latest restore $go_encoded)
    python_decoded=$(echo $python_decoded | sed 's/^"//;s/"$//')
    echo "Python Decoded: $python_decoded"

    # Check if the Python decoded matches the original secret
    if [ "$python_decoded" != "$secret" ]; then
        echo "Test failed: Python failed to decode Go encoded secret"
        exit 1
    fi

    echo "All tests passed successfully!"
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
        IFS=' ' read -r -a threshold_array <<< "$threshold"
        run_test "$secret" "${threshold_array[0]}" "${threshold_array[1]}"
        result=$(echo $?)
        if [ "$result" != "0" ]; then
            echo "Test failed. Stopping: $result"
            exit 1
        fi
        echo "Test passed for secret: '$secret' with threshold: $threshold"
        echo "----------------------------------------"
    done
done