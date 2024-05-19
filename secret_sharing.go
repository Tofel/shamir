// package main

// import (
// 	"crypto/rand"
// 	"encoding/hex"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// // PRIME is a large prime number used for the finite field
// var PRIME = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))

// // modInverse calculates the modular multiplicative inverse of a modulo prime
// func modInverse(a, prime *big.Int) *big.Int {
// 	return new(big.Int).ModInverse(a, prime)
// }

// // evalAt evaluates the polynomial at a given x using Horner's method
// func evalAt(poly []*big.Int, x, prime *big.Int) *big.Int {
// 	accum := new(big.Int)
// 	for i := len(poly) - 1; i >= 0; i-- {
// 		accum.Mul(accum, x)
// 		accum.Add(accum, poly[i])
// 		accum.Mod(accum, prime)
// 	}
// 	return accum
// }

// // makeShares creates the shares using Shamir's Secret Sharing
// func makeShares(secret *big.Int, numShares, threshold int, prime *big.Int) ([]*big.Int, error) {
// 	if threshold > numShares {
// 		return nil, errors.New("threshold cannot be greater than number of shares")
// 	}

// 	// Generate a random polynomial with the given secret as the intercept
// 	poly := make([]*big.Int, threshold)
// 	poly[0] = secret
// 	for i := 1; i < threshold; i++ {
// 		randCoeff, _ := rand.Int(rand.Reader, prime)
// 		poly[i] = randCoeff
// 	}

// 	// Generate the shares
// 	shares := make([]*big.Int, numShares)
// 	for i := 1; i <= numShares; i++ {
// 		shares[i-1] = evalAt(poly, big.NewInt(int64(i)), prime)
// 	}

// 	return shares, nil
// }

// // lagrangeInterpolate performs Lagrange interpolation to recover the secret
// func lagrangeInterpolate(x *big.Int, xVals, yVals []*big.Int, prime *big.Int) *big.Int {
// 	k := len(xVals)
// 	accum := new(big.Int)
// 	PI := func(vals []*big.Int) *big.Int {
// 		result := big.NewInt(1)
// 		for _, v := range vals {
// 			result.Mul(result, v)
// 		}
// 		return result
// 	}

// 	for i := 0; i < k; i++ {
// 		numer := PI(removeIndex(xVals, i))
// 		denom := PI(subtractEach(xVals[i], removeIndex(xVals, i)))
// 		term := new(big.Int).Mul(yVals[i], numer)
// 		term.Mul(term, modInverse(denom, prime))
// 		accum.Add(accum, term)
// 	}

// 	accum.Mod(accum, prime)
// 	return accum
// }

// // removeIndex removes an element from a slice at a given index
// func removeIndex(slice []*big.Int, index int) []*big.Int {
// 	result := make([]*big.Int, len(slice)-1)
// 	copy(result[0:], slice[0:index])
// 	copy(result[index:], slice[index+1:])
// 	return result
// }

// // subtractEach subtracts a value from each element in a slice
// func subtractEach(val *big.Int, slice []*big.Int) []*big.Int {
// 	result := make([]*big.Int, len(slice))
// 	for i, v := range slice {
// 		result[i] = new(big.Int).Sub(val, v)
// 	}
// 	return result
// }

// // splitSecret splits a secret into shares
// func splitSecret(secret []byte, numShares, threshold int) ([]string, error) {
// 	secretInt := new(big.Int).SetBytes(secret)
// 	shares, err := makeShares(secretInt, numShares, threshold, PRIME)
// 	if err != nil {
// 		return nil, err
// 	}

// 	maxLength := 0
// 	for _, share := range shares {
// 		if len(share.Bytes()) > maxLength {
// 			maxLength = len(share.Bytes())
// 		}
// 	}

// 	shareStrs := make([]string, numShares)
// 	for i, share := range shares {
// 		paddedShare := fmt.Sprintf("%0*x", maxLength*2, share)
// 		shareStrs[i] = fmt.Sprintf("%d-%s", i+1, paddedShare)
// 	}

// 	// Assert that the number of shares matches the expected number
// 	if len(shareStrs) != numShares {
// 		return nil, errors.New("number of shares does not match expected number")
// 	}

// 	// Assert that the shares are correctly formatted
// 	for _, share := range shareStrs {
// 		if !strings.Contains(share, "-") {
// 			return nil, errors.New("share format is incorrect")
// 		}
// 	}

// 	return shareStrs, nil
// }

// // recoverSecret recovers the secret from shares
// func recoverSecret(shareStrs []string) ([]byte, error) {
// 	xVals := make([]*big.Int, len(shareStrs))
// 	yVals := make([]*big.Int, len(shareStrs))

// 	for i, shareStr := range shareStrs {
// 		parts := strings.Split(shareStr, "-")
// 		if len(parts) != 2 {
// 			return nil, errors.New("invalid share format")
// 		}

// 		x, err := strconv.Atoi(parts[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		y, err := hex.DecodeString(parts[1])
// 		if err != nil {
// 			return nil, err
// 		}

// 		xVals[i] = big.NewInt(int64(x))
// 		yVals[i] = new(big.Int).SetBytes(y)
// 	}

// 	// Assert that the number of x values matches the number of y values
// 	if len(xVals) != len(yVals) {
// 		return nil, errors.New("number of x values does not match number of y values")
// 	}

// 	secretInt := lagrangeInterpolate(big.NewInt(0), xVals, yVals, PRIME)
// 	secretBytes := secretInt.Bytes()

// 	// Assert that the recovered secret is not empty
// 	if len(secretBytes) == 0 {
// 		return nil, errors.New("recovered secret is empty")
// 	}

// 	return secretBytes, nil
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: main <split|restore> [options]")
// 		os.Exit(1)
// 	}

// 	action := os.Args[1]

// 	switch action {
// 	case "split":
// 		if len(os.Args) != 5 {
// 			fmt.Println("Usage: main split <input_string> <num_shares> <threshold>")
// 			os.Exit(1)
// 		}

// 		inputString := os.Args[2]
// 		numShares, err := strconv.Atoi(os.Args[3])
// 		if err != nil {
// 			fmt.Printf("Invalid number of shares: %v\n", err)
// 			os.Exit(1)
// 		}

// 		threshold, err := strconv.Atoi(os.Args[4])
// 		if err != nil {
// 			fmt.Printf("Invalid threshold: %v\n", err)
// 			os.Exit(1)
// 		}

// 		shares, err := splitSecret([]byte(inputString), numShares, threshold)
// 		if err != nil {
// 			fmt.Printf("Error splitting string: %v\n", err)
// 			os.Exit(1)
// 		}

// 		fmt.Println("Shares:")
// 		for _, share := range shares {
// 			fmt.Println(share)
// 		}

// 	case "restore":
// 		if len(os.Args) != 3 {
// 			fmt.Println("Usage: main restore <comma_separated_shares>")
// 			os.Exit(1)
// 		}

// 		shareStrs := strings.Split(os.Args[2], ",")
// 		restoredSecret, err := recoverSecret(shareStrs)
// 		if err != nil {
// 			fmt.Printf("Error restoring secret: %v\n", err)
// 			os.Exit(1)
// 		}

// 		fmt.Println("Restored Secret:")
// 		fmt.Println(string(restoredSecret))

// 	default:
// 		fmt.Println("Unknown action:", action)
// 		fmt.Println("Usage: main <split|restore> [options]")
// 		os.Exit(1)
// 	}
// }
