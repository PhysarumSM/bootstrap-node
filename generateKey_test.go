package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateKey(test *testing.T) {
	// Setup for ExistngKey case
	tmpFile, err := ioutil.TempFile("/tmp", "tmp")
	if err != nil {
		panic(err)
	}
	test.Log("Pre-created temp key file:", tmpFile.Name())
	tmpFile.Close()

	testCases := []struct {
		name      string
		algo      string
		bits      int
		keyFile   string
		shouldErr bool
	}{
		// Negative test cases
		{"ExistingKey", "rsa", 2048, tmpFile.Name(), true},
		{"RSA-small-bits", "rsa", 1024, "~/asdf", true},
		{"Wrong-algo", "rsaa", 2048, "~/asdf", true},

		// Positive test cases
		{"RSA-basic", "rsa", 2048, "~/asdf", false},
		{"ECDSA-basic", "ecdsa", 2048, "~/asdf", false},
		{"Ed25519-basic", "ecdsa", 0, "~/asdf", false},
		{"Ed25519-basic", "Ed25519", 0, "~/asdf", false},
		{"Secp256k1-basic", "Secp256k1", 0, "~/asdf", false},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			err := generateKey(testCase.algo, testCase.bits, testCase.keyFile)
			if testCase.shouldErr {
				if err == nil {
					test.Errorf("Passed case (%s); Expected it to fail.", testCase.name)
				} else {
					test.Log(err)
				}
			} else if !testCase.shouldErr && err != nil {
				test.Log(err)
				test.Errorf("Failed case (%s); Expected it to pass.\n%v", testCase.name, err)
			}

			// This assumes expandTilde was fully tested...
			if keyFilePath, err := expandTilde(testCase.keyFile); err == nil {
				os.Remove(keyFilePath)
			}
		})
	}

}
