package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
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
)

func main() {
	flag.Parse()

	if *genKey {
		fmt.Println("Generating a new key...")
		priv, err := generatePrivKey(*algo, *bits)
		if err != nil {
			fmt.Printf("ERROR: Unable to generate key\n%v", err)
			os.Exit(1)
		}

		if err = storePrivKeyToFile(priv, *keyFile); err != nil {
			fmt.Printf("ERROR: Unable to save key to file %s\n", *keyFile)
			os.Exit(1)
		}
		fmt.Println("New key is stored at:", *keyFile)
	}

	if !fileExists(*keyFile) {
		fmt.Printf("ERROR: Key (%s) does not exist.\n", *keyFile)
		fmt.Printf("Ensure path is correct or generate a new key with -genkey.\n")
		os.Exit(1)
	}

	priv, err := loadPrivKeyFromFile(*keyFile)
	if err != nil {
		fmt.Printf("ERROR: Unable to load key from file\n%v", err)
		os.Exit(1)
	}

	listenAddrs, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
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
