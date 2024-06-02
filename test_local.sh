#!/bin/bash

shuffle() {
   local i tmp size max rand

   # $RANDOM % (i+1) is biased because of the limited range of $RANDOM
   # Compensate by using a range which is a multiple of the array size.
   size=${#array[*]}
   max=$(( 32768 / size * size ))

   for ((i=size-1; i>0; i--)); do
      while (( (rand=$RANDOM) >= max )); do :; done
      rand=$(( rand % (i+1) ))
      tmp=${array[i]} array[i]=${array[rand]} array[rand]=$tmp
   done
}

# Function to run encoding and decoding tests
run_test() {
    secret="$1"
    threshold="$2"
    total_shares="$3"

    # Python encode
    python_encoded=$(python3 shamir.py split "$secret" "$threshold" "$total_shares")
    # echo "Python Encoded: $python_encoded"

    # Go decode Python encoded
    go_decoded=$(go run shamir.go restore "$python_encoded")
    # echo "Go Decoded: $go_decoded"

    # Check if the Go decoded matches the original secret
    if [ "$go_decoded" != "$secret" ]; then
        echo "Test failed: Go failed to decode Python encoded secret"
        exit 1
    fi

    # Go encode
    go_encoded=$(go run shamir.go split "$secret" "$threshold" "$total_shares")
    # echo "Go Encoded: $go_encoded"

    # Python decode Go encoded
    python_decoded=$(python3 shamir.py restore "$go_encoded")
    # echo "Python Decoded: $python_decoded"

    # Check if the Python decoded matches the original secret
    if [ "$python_decoded" != "$secret" ]; then
        echo "Test failed: Python failed to decode Go encoded secret"
        exit 1
    fi

    # Convert comma-separated python_encoded string into an array
    IFS=',' read -r -a shares_array <<< "$python_encoded"

    # Iterate over the original array and construct inputs with various numbers of elements
    for (( n=1; n<=${#shares_array[@]}; n++ )); do
        echo -n "Testing with $n/$total_shares ordered shares with $threshold threshold..."
        test_array=()
        for (( i=0; i<$n; i++ )); do
            test_array+=("${shares_array[i]}")
        done

        shares_input=$(IFS=,; echo "${test_array[*]}")

        # Decode with n shares
        if [ "$n" -lt "$threshold" ]; then
            # Expect an error
            go_decoded=$(go run shamir.go restore "$shares_input" 2>&1)
            if [ "$go_decoded" = "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares ordered shares with $threshold threshold... NOK! --> Expected error when decoding with Go with $n shares, but got success"
                exit 1
            fi
        else
            # Expect success
            go_decoded=$(go run shamir.go restore "$shares_input")
            if [ "$go_decoded" != "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares ordered shares with $threshold threshold... NOK! --> Go failed to decode Python encoded secret with $n shares"
                exit 1
            fi
        fi

        # Repeat the process for Go encoded shares
        if [ "$n" -lt "$threshold" ]; then
            # Expect an error
            python_decoded=$(python3 shamir.py restore "$shares_input" 2>&1)
            if [ "$python_decoded" = "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares ordered shares with $threshold threshold... NOK! --> Expected error when decoding with Python with $n shares, but got success"
                exit 1
            fi
        else
            # Expect success
            python_decoded=$(python3 shamir.py restore "$shares_input")
            if [ "$python_decoded" != "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares ordered shares with $threshold threshold... NOK! --> Python failed to decode Go encoded secret with $n shares"
                exit 1
            fi
        fi

        echo -ne "\rTesting with $n/$total_shares ordered shares with $threshold threshold... OK!" 
        echo
        echo -n "Testing with $n/$total_shares shuffled shares with $threshold threshold..."
        array=("${test_array[@]}")
        shuffle

        shares_input=$(IFS=,; echo "${array[*]}")

        # Decode with n shares
        if [ "$n" -lt "$threshold" ]; then
            # Expect an error
            go_decoded=$(go run shamir.go restore "$shares_input" 2>&1)
            if [ "$go_decoded" = "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares shuffled shares with $threshold threshold... NOK! --> Expected error when decoding with Go with $n shares, but got success"
                exit 1
            fi
        else
            # Expect success
            go_decoded=$(go run shamir.go restore "$shares_input")
            if [ "$go_decoded" != "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares shuffled shares with $threshold threshold... NOK! --> Go failed to decode Python encoded secret with $n shares"
                exit 1
            fi
            # echo
            # echo "secret decoded with Go from shuffled shares: $go_decoded"
        fi

        # Repeat the process for Go encoded shares
        if [ "$n" -lt "$threshold" ]; then
            # Expect an error
            python_decoded=$(python3 shamir.py restore "$shares_input" 2>&1)
            if [ "$python_decoded" = "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares shuffled shares with $threshold threshold... NOK! --> Expected error when decoding with Python with $n shares, but got success"
                exit 1
            fi
        else
            # Expect success
            python_decoded=$(python3 shamir.py restore "$shares_input")
            if [ "$python_decoded" != "$secret" ]; then
                echo -ne "\rTesting with $n/$total_shares shuffled shares with $threshold threshold... NOK! --> Python failed to decode Go encoded secret with $n shares"
                exit 1
            fi
            # echo "secret decoded with Python from shuffled shares: $python_decoded"
        fi        

        echo -ne "\rTesting with $n/$total_shares shuffled shares with $threshold threshold... OK!" 
        echo
        echo "Tests passed for $n/$total_shares shares"
        echo "----------------------------------------"
    done

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
        echo "-----> Testing secret: '$secret' with threshold: $threshold"
        echo
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