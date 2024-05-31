package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

func splitSecret(secret string, totalShares int, threshold int) (string, error) {
	shares, err := shamir.Split([]byte(secret), totalShares, threshold)
	if err != nil {
		return "", err
	}

	encodedShares := make([]string, totalShares)
	for i, share := range shares {
		idx := i
		encodedShares[i] = fmt.Sprintf("%d-%s", idx+1, hex.EncodeToString(share))
	}
	return strings.Join(encodedShares, ","), nil
}

func restoreSecret(encodedShares string) (string, error) {
	shareStrings := strings.Split(encodedShares, ",")
	shares := make([][]byte, len(shareStrings))

	for i, shareStr := range shareStrings {
		parts := strings.SplitN(shareStr, "-", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid share format")
		}
		share, err := hex.DecodeString(parts[1])
		if err != nil {
			return "", err
		}
		shares[i] = share
	}

	secret, err := shamir.Combine(shares)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run shamir.go <command> <args>")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "split":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run shamir.go split <secret> <threshold> <total_shares>")
			os.Exit(1)
		}
		secret := os.Args[2]
		threshold := os.Args[3]
		totalShares := os.Args[4]

		totalSharesInt, err := strconv.Atoi(totalShares)
		if err != nil {
			fmt.Println("Invalid total_shares value")
			os.Exit(1)
		}

		thresholdInt, err := strconv.Atoi(threshold)
		if err != nil {
			fmt.Println("Invalid threshold value")
			os.Exit(1)
		}

		if thresholdInt > totalSharesInt {
			fmt.Println("Threshold cannot be bigger than total shares")
			os.Exit(1)
		}

		encoded, err := splitSecret(secret, totalSharesInt, thresholdInt)
		if err != nil {
			fmt.Printf("Error splitting secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(encoded)

	case "restore":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run shamir.go restore <encoded_shares>")
			os.Exit(1)
		}
		encodedShares := os.Args[2]

		secret, err := restoreSecret(encodedShares)
		if err != nil {
			fmt.Printf("Error restoring secret: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(secret)

	default:
		fmt.Println("Invalid command. Use 'split' or 'restore'")
		os.Exit(1)
	}
}
