package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	RSA_MIN_BITS = 2048
)

var (
	// Map algos to enum from libp2p-core/crypto/key
	// See: https://github.com/libp2p/go-libp2p-core/blob/master/crypto/key.go
	// TODO: Figure out future-proof way to maintain this? i.e. avoid manual
	//       modification here if libp2p adds new algos.
	keyTypes = map[string]int{
		"rsa":       crypto.RSA,
		"ed25519":   crypto.Ed25519,
		"secp256k1": crypto.Secp256k1,
		"ecdsa":     crypto.ECDSA,
	}
)

// Expands tilde to absolute path
// Currently only works if path begins with tilde, not somewhere in the middle
func expandTilde(path string) (string, error) {
	newPath := path

	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err != nil {
			return "", err
		} else {
			newPath = home + path[1:]
		}
	}

	return newPath, nil
}

func generateKey(algo string, bits int, keyFile string) error {
	var keyType int
	for algoName, algoID := range keyTypes {
		if strings.EqualFold(algoName, algo) {
			keyType = algoID
			break
		}
		keyType = -1
	}

	if keyType < 0 {
		return fmt.Errorf("Unknown algorithm")
	} else if keyType == crypto.RSA && bits < RSA_MIN_BITS {
		return fmt.Errorf("Number of bits for RSA must be at least %d", RSA_MIN_BITS)
	}

	// Generate key, then write to file
	priv, _, err := crypto.GenerateKeyPair(keyType, bits)
	if err != nil {
		return err
	}

	if keyFile, err = expandTilde(keyFile); err != nil {
		return err
	}

	_, err = os.Stat(keyFile)
	if !os.IsNotExist(err) {
		return fmt.Errorf("A key file already exists (%s).\n"+
			"Delete it or move it before proceeding.", keyFile)
	}

	file, err := os.Create(keyFile)
	if err != nil {
		return err
	} else {
		defer file.Close()
	}

	rawBytes, err := priv.Raw()
	if err != nil {
		return err
	}
	_, err = file.WriteString(crypto.ConfigEncodeKey(rawBytes))
	if err != nil {
		return err
	}

	return nil
}
