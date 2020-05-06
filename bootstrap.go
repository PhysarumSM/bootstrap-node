package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

var (
	genKey = flag.Bool("genkey", false,
		"Generate a new key and save to file.")
	algo = flag.String("algo", "RSA",
		"Cryptographic algorithm to use for generating the key.\n"+
			"Will be ignored if 'gen-key' is false.\n"+
			"Must be one of {RSA, Ed25519, Secp256k1, ECDSA}")
	bits = flag.Int("bits", 2048,
		"Key length, in bits. Will be ignored if 'algo' is not RSA.")
	keyFile = flag.String("keyfile", "~/.privKey",
		"Location of private key to read from (or write to, if generating).")

	// Map algos to enum from libp2p-core/crypto/key
	// See: https://github.com/libp2p/go-libp2p-core/blob/master/crypto/key.go
	// TODO: Figure out future-proof way to maintain this?
	keyTypes = map[string]int{
		"rsa":       crypto.RSA,
		"ed25519":   crypto.Ed25519,
		"secp256k1": crypto.Secp256k1,
		"ecdsa":     crypto.ECDSA,
	}
)

func main() {
	flag.Parse()

	// TODO: Clean up this init section
	//       Move to function + cleaner error msg and exit mechanism
	if *genKey {
		var keyType int
		for algoName, algoID := range keyTypes {
			if strings.EqualFold(algoName, *algo) {
				keyType = algoID
				break
			}
			keyType = -1
		}

		if keyType < 0 {
			fmt.Println("ERROR: Unknown algorithm")
			os.Exit(1)
		} else if keyType == crypto.RSA && *bits < 1024 {
			fmt.Println("ERROR: Number of bits for RSA must be at least 1024")
			os.Exit(1)
		}

		// Generate key, then write to file
		priv, _, err := crypto.GenerateKeyPair(keyType, *bits)
		if err != nil {
			panic(err)
		}

        if strings.HasPrefix(*keyFile, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			*keyFile = home + (*keyFile)[1:]
		}

		_, err = os.Stat(*keyFile)
		if !os.IsNotExist(err) {
			fmt.Println("ERROR: A key file already exists at", *keyFile)
			fmt.Println("       Delete it or move it before proceeding")
			os.Exit(1)
		}

		file, err := os.Create(*keyFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		rawBytes, err := priv.Raw()
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(crypto.ConfigEncodeKey(rawBytes))
		if err != nil {
			panic(err)
		}
	} else {
		// TODO: Read key from file
		fmt.Println("read from file...")
	}

	ctx := context.Background()
	// Set your own keypair
	// priv, _, err := crypto.GenerateKeyPair(
	// 	crypto.Ed25519, // Select your key type. Ed25519 are nice short
	// 	-1,             // Select key length when possible (i.e. RSA).
	// )
	// priv, err := rsa.GenerateKey(rand.Reader, bits)
	// fmt.Println("Private key INFO:", priv)
	// var DefaultPeerstore Option = func(cfg *Config) error {
	//     return cfg.Apply(Peerstore(pstoremem.NewPeerstore("QmZJexJchxXt71N9bSkj6PchxyGEs9g6Qvrj7fEk32FqDs")))
	// }

	byteVal := []byte{211, 151, 127, 224, 159, 14, 157, 18, 23, 132, 211, 171, 4, 8, 125, 131, 235, 83, 169, 205, 79, 230, 32, 138, 150, 179, 103, 28, 152, 240, 11, 111, 101, 134, 246, 174, 231, 186, 183, 172, 59, 180, 89, 156, 126, 43, 240, 153, 190, 62, 31, 24, 209, 96, 245, 188, 19, 240, 39, 95, 93, 41, 140, 38}
	priv, err := crypto.UnmarshalEd25519PrivateKey(byteVal)
	if err != nil {
		panic(err)
	}
	listenAddrs, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err != nil {
		panic(err)
	}
	node, err := libp2p.New(ctx, libp2p.ListenAddrs(listenAddrs), libp2p.Identity(priv))
	if err != nil {
		panic(err)
	}
	fmt.Println("This node: ", node.ID().Pretty(), " ", node.Addrs())
	_, err = dht.New(ctx, node)
	if err != nil {
		panic(err)
	}

	select {}
}
