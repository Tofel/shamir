package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Define a known large prime number for the finite field
var PRIME, _ = new(big.Int).SetString("134078079299425970995740249982058461273830392263134615011048558635531023870321", 10)

// Generate random coefficients for the polynomial
func generateCoefficients(threshold int) ([]*big.Int, error) {
	coefficients := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		coeff, err := rand.Int(rand.Reader, PRIME)
		if err != nil {
			return nil, err
		}
		coefficients[i] = coeff
	}
	return coefficients, nil
}

// Evaluate the polynomial at a given x
func evalAt(poly []*big.Int, x *big.Int, prime *big.Int) *big.Int {
	result := new(big.Int).Set(poly[0])
	xPower := new(big.Int).Set(x)
	for i := 1; i < len(poly); i++ {
		term := new(big.Int).Mul(poly[i], xPower)
		result.Add(result, term)
		result.Mod(result, prime)
		xPower.Mul(xPower, x)
		xPower.Mod(xPower, prime)
	}
	return result
}

// Pad the share to ensure consistent length
func padShare(share *big.Int, length int) string {
	shareBytes := share.Bytes()
	if len(shareBytes) < length {
		padding := make([]byte, length-len(shareBytes))
		shareBytes = append(padding, shareBytes...)
	}
	return hex.EncodeToString(shareBytes)
}

// Split the secret into shares
func splitSecret(secret []byte, numShares, threshold int) ([]string, error) {
	secretInt := new(big.Int).SetBytes(secret)
	coefficients, err := generateCoefficients(threshold - 1)
	if err != nil {
		return nil, err
	}
	coefficients = append([]*big.Int{secretInt}, coefficients...) // Prepend the secret as the first coefficient

	// Determine the maximum byte length for padding
	maxByteLength := len(secretInt.Bytes())
	for _, coeff := range coefficients {
		if len(coeff.Bytes()) > maxByteLength {
			maxByteLength = len(coeff.Bytes())
		}
	}

	shares := make([]string, numShares)
	for i := 1; i <= numShares; i++ {
		x := big.NewInt(int64(i))
		y := evalAt(coefficients, x, PRIME)
		paddedShare := padShare(y, maxByteLength)
		share := fmt.Sprintf("%d-%s", i, paddedShare)
		shares[i-1] = share
	}
	return shares, nil
}

// Lagrange interpolation to reconstruct the secret
func lagrangeInterpolate(x, prime *big.Int, xVals, yVals []*big.Int) (*big.Int, error) {
	result := new(big.Int)
	k := len(xVals)
	for i := 0; i < k; i++ {
		num := new(big.Int).SetInt64(1)
		den := new(big.Int).SetInt64(1)
		for j := 0; j < k; j++ {
			if i == j {
				continue
			}
			num.Mul(num, new(big.Int).Neg(xVals[j]))
			num.Mod(num, prime)
			den.Mul(den, new(big.Int).Sub(xVals[i], xVals[j]))
			den.Mod(den, prime)
		}
		fmt.Printf("xVals[%d] = %s, num = %s, den = %s\n", i, xVals[i].String(), num.String(), den.String())
		if den.Sign() == 0 {
			return nil, fmt.Errorf("denominator is zero for index %d", i)
		}
		inv := new(big.Int).ModInverse(den, prime)
		if inv == nil {
			return nil, fmt.Errorf("modular inverse does not exist for den %s", den.String())
		}
		term := new(big.Int).Mul(yVals[i], num)
		term.Mul(term, inv)
		term.Mod(term, prime)
		result.Add(result, term)
		result.Mod(result, prime)
	}
	return result, nil
}

// Recover the secret from shares
func recoverSecret(shareStrs []string) ([]byte, error) {
	xVals := make([]*big.Int, len(shareStrs))
	yVals := make([]*big.Int, len(shareStrs))

	for i, shareStr := range shareStrs {
		parts := strings.Split(shareStr, "-")
		if len(parts) != 2 {
			return nil, errors.New("invalid share format")
		}

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		yPadded, err := hex.DecodeString(parts[1])
		if err != nil {
			return nil, err
		}

		xVals[i] = big.NewInt(int64(x))
		yVals[i] = new(big.Int).SetBytes(yPadded)
	}

	secretInt, err := lagrangeInterpolate(big.NewInt(0), PRIME, xVals, yVals)
	if err != nil {
		return nil, err
	}
	secretBytes := secretInt.Bytes()
	return secretBytes, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <split|restore> [options]")
		os.Exit(1)
	}

	action := os.Args[1]

	switch action {
	case "split":
		if len(os.Args) != 5 {
			fmt.Println("Usage: main split <input_string> <num_shares> <threshold>")
			os.Exit(1)
		}

		inputString := os.Args[2]
		numShares, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("Invalid number of shares: %v\n", err)
			os.Exit(1)
		}

		threshold, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Printf("Invalid threshold: %v\n", err)
			os.Exit(1)
		}

		shares, err := splitSecret([]byte(inputString), numShares, threshold)
		if err != nil {
			fmt.Printf("Error splitting string: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Shares:")
		for _, share := range shares {
			fmt.Println(share)
		}

	case "restore":
		if len(os.Args) != 3 {
			fmt.Println("Usage: main restore <comma_separated_shares>")
			os.Exit(1)
		}

		shareStrs := strings.Split(os.Args[2], ",")
		restoredSecret, err := recoverSecret(shareStrs)
		if err != nil {
			fmt.Printf("Error restoring secret: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Restored Secret:")
		fmt.Println(string(restoredSecret))

	default:
		fmt.Println("Unknown action:", action)
		fmt.Println("Usage: main <split|restore> [options]")
		os.Exit(1)
	}
}
