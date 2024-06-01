package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitAndRestoreSecret(t *testing.T) {
	type testCase struct {
		name        string
		secret      string
		threshold   int
		totalShares int
	}

	testCases := []testCase{
		{
			name:        "simple word",
			secret:      "my_secret",
			threshold:   2,
			totalShares: 3,
		},
		{
			name:        "simple sentece",
			secret:      "god of mices ate my food and then he cried",
			threshold:   2,
			totalShares: 6,
		},
		{
			name:        "long sentece",
			secret:      "god of mices ate my food and then he cried while other mices danced around him",
			threshold:   6,
			totalShares: 6,
		},
		{
			name:        "BIP-39",
			secret:      "pen aunt text rotate donate sock shield pottery cloud toy tank sibling parrot oblige agent egg october angle short wolf survey frequent autumn desert",
			threshold:   3,
			totalShares: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secret := tc.secret
			threshold := tc.threshold
			totalShares := tc.totalShares

			encodedShares, err := splitSecret(secret, totalShares, threshold)
			require.NoError(t, err, "splitting should not fail")

			restoredSecret, err := restoreSecret(encodedShares)
			require.NoError(t, err, "restoring secret from all shares should not fail")
			require.Equal(t, secret, restoredSecret, "restored secret should be the same as original one")

			shuffledShares := shuffleShares(encodedShares)
			require.NotEqual(t, encodedShares, shuffledShares, "shuffled shares should be different from original ones")

			restoredSecret, err = restoreSecret(shuffledShares)
			require.NoError(t, err, "restoring secret from shuffled shares should not fail")
			require.Equal(t, secret, restoredSecret, "restored secret should be the same as original one")

			for i := 1; i < totalShares; i++ {
				if i < threshold {
					insufficientShares := getNshares(encodedShares, i)
					fmt.Printf("[%d]insufficientShares: %s\n", i, insufficientShares)
					restoredSecret, err := restoreSecret(insufficientShares)
					// quirk of the library, insufficient shares > 1 are decoded, but secret is incorrect
					if i == 1 {
						require.Error(t, err, "restoring secret with 1 share should fail")
					} else {
						require.NotEqual(t, secret, restoredSecret, "restoring secret with insufficient shares should result in incorrect secret")
					}
				} else {
					suffcientShares := getNshares(encodedShares, i)
					restoredSecret, err = restoreSecret(suffcientShares)
					require.NoError(t, err, "restoring shares with shares >= threshold should not fail")
					require.Equal(t, secret, restoredSecret, "restored secret should be the same as original one")
				}
			}

			incompleteShares := encodedShares[:len(encodedShares)-2] // Remove part of last share
			_, err = restoreSecret(incompleteShares)
			require.Error(t, err, "restoring secret with trimmed shares should fail")

			extraShares := addExtraShares(encodedShares)
			_, err = restoreSecret(extraShares)
			require.Error(t, err, "restoring secret with incorrect extra share should fail")

			_, err = restoreSecret("1624ghjsgd762")
			require.Error(t, err, "restoring secret from malformed data should fail")

			_, err = restoreSecret(removeIndexesFromShares(encodedShares))
			require.Error(t, err, "restoring secret from shares without indexes should fail")
		})
	}
}

func TestIncorrectSplit(t *testing.T) {
	type testCase struct {
		name        string
		threshold   int
		totalShares int
	}

	testCases := []testCase{
		{
			name:        "zero threshold",
			threshold:   0,
			totalShares: 3,
		},
		{
			name:        "too low threshold",
			threshold:   1,
			totalShares: 3,
		},
		{
			name:        "threshold > total shares",
			threshold:   3,
			totalShares: 2,
		},
		{
			name:        "both zero",
			threshold:   0,
			totalShares: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			threshold := tc.threshold
			totalShares := tc.totalShares

			_, err := splitSecret("my_secret", totalShares, threshold)
			require.Error(t, err, "splitting with incorrect inputs should fail")
		})
	}
}

func shuffleShares(encodedShares string) string {
	shares := strings.Split(encodedShares, ",")
	rand.Shuffle(len(shares), func(i, j int) {
		shares[i], shares[j] = shares[j], shares[i]
	})
	return strings.Join(shares, ",")
}

func getNshares(encodedShares string, newCount int) string {
	shares := strings.Split(encodedShares, ",")
	var newShares []string
	for i := 0; i < newCount; i++ {
		newShares = append(newShares, shares[i])
	}
	return strings.Join(newShares, ",")
}

func removeIndexesFromShares(encodedShares string) string {
	shares := strings.Split(encodedShares, ",")
	var newShares []string
	for _, share := range shares {
		split := strings.Split(share, "-")
		newShares = append(newShares, split[1])
	}
	return strings.Join(newShares, ",")
}

func addExtraShares(encodedShares string) string {
	extraShare := "4-" + hex.EncodeToString([]byte("extra_share_data"))
	return encodedShares + "," + extraShare
}
