package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
)

func TestLoadKey(test *testing.T) {
	// Create an existing key to load from
	keyType := pb.KeyType(3)
	keyB64 := "MHcCAQEEIHp/bhcT3Jge9ykOMjk+AgCi6qqM8it01IRoRbXphHXaoAoGCCqGSM49AwEHoUQDQgAEhN7JYn9DN9POlfbkDwR1T74gxPpUx90cWxbuyuvOL10DsQe1UD/IVBxdQ1nZPaYC/m+nSaUdZ53gFBaHLQg+QQ=="

	tmpFile, err := ioutil.TempFile("/tmp", "tmp")
	if err != nil {
		panic(err)
	}
	tmpFile.WriteString(fmt.Sprintf("%d %s\n", keyType, keyB64))
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	priv, err := loadPrivKeyFromFile(tmpFile.Name())
	if err != nil {
		test.Fatalf("loadPrivKeyFromFile() failed with error:\n%v", err)
	}

	if priv.Type() != keyType {
		test.Fatalf("Incorrect key type loaded (%d), was expecting %d", priv.Type(), keyType)
	}

	rawBytes, err := priv.Raw()
	if err != nil {
		test.Fatalf("Could not load raw bytes from loaded key")
	}

	loadedKeyB64 := crypto.ConfigEncodeKey(rawBytes)
	if loadedKeyB64 != keyB64 {
		test.Fatalf("Loaded key is not identical to test key.\n"+
			"Loaded: %s\nExpect: %s\n", loadedKeyB64, keyB64)
	}
}
