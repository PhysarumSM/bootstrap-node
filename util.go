package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
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

func fileExists(filePath string) bool {
	filePath, err := expandTilde(filePath)
	if err != nil {
		return false
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}

	return true
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

	// TODO: Ideally this next part should be moved out of the function.
	//       It'll make the function more flexible/reusable, as users may
	//       want to generate a key but not store it.
	if err = storePrivKeyToFile(priv, keyFile); err != nil {
		return err
	}

	return nil
}

// Write private key to file in Base 64 format
// Store the key type ID followed by a space, then the key, then a new-line
func storePrivKeyToFile(priv crypto.PrivKey, keyFile string) error {
	keyFile, err := expandTilde(keyFile)
	if err != nil {
		return err
	}

	if fileExists(keyFile) {
		return fmt.Errorf("File already exists (%s).\n"+
			"Delete it or move it before proceeding.", keyFile)
	}

	file, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer file.Close()

	rawBytes, err := priv.Raw()
	if err != nil {
		return err
	}

	fileStr := fmt.Sprintf("%d %s\n", priv.Type(), crypto.ConfigEncodeKey(rawBytes))
	_, err = file.WriteString(fileStr)
	if err != nil {
		return err
	}

	return nil
}

// Inverse of storePrivKeyToFile()
func loadPrivKeyFromFile(keyFile string) (crypto.PrivKey, error) {
	keyFile, err := expandTilde(keyFile)
	if err != nil {
		return nil, err
	}

	if !fileExists(keyFile) {
		return nil, fmt.Errorf("File (%s) does not exist.", keyFile)
	}

	/* NOTE: Using ioutil's ReadFile() may be potentially bad in the case that
	 *       the file is very large, as it tries to read the entire file at once.
	 *       Alternative is to read chunk by chunk using os' Read() and combine.
	 *       I'm being lazy, assume file is small or memory is large.
	 */
	content, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	// Strip new-line, then parse key type from key itself
	contentStr := string(content[:len(content)-1])
	spaceIdx := strings.IndexByte(contentStr, ' ')
	if spaceIdx <= 0 {
		return nil, fmt.Errorf("Unable to load key file (may have been corrupted)")
	}

	keyType, err := strconv.ParseInt(contentStr[:spaceIdx], 10, 32)
	if err != nil {
		return nil, err
	}

	keyB64 := contentStr[spaceIdx+1:]
	keyRaw, err := crypto.ConfigDecodeKey(keyB64)
	if err != nil {
		return nil, err
	}

	// Unmarsall to create private key object
	unmarshaller, ok := crypto.PrivKeyUnmarshallers[pb.KeyType(keyType)]
	if !ok {
		return nil, fmt.Errorf("Key file contains an unknown algorithm.")
	}

	return unmarshaller(keyRaw)
}
