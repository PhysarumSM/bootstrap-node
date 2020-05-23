/* Copyright 2020 Multi-Tier-Cloud Development Team
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/Multi-Tier-Cloud/common/p2pnode"
	"github.com/Multi-Tier-Cloud/common/util"
)

var (
	genKey = flag.Bool("genkey", false,
		"Generate a new key and save to file.")
	algo = flag.String("algo", "RSA",
		"Cryptographic algorithm to use for generating the key.\n"+
			"Will be ignored if 'genkey' is false.\n"+
			"Must be one of {RSA, Ed25519, Secp256k1, ECDSA}")
	bits = flag.Int("bits", 2048,
		"Key length, in bits. Will be ignored if 'algo' is not RSA.")
	keyFile = flag.String("keyfile", "~/.privKey",
		"Location of private key to read from (or write to, if generating).")
	ephemeral = flag.Bool("ephemeral", false,
		"Generate a new key just for this run, and don't store it to file.\n"+
			"If 'keyfile' is specified, it will be ignored.")
)

func main() {
	flag.Parse()

	var priv crypto.PrivKey
	var err error
	if *genKey || *ephemeral {
		fmt.Println("Generating a new key...")
		priv, err = util.GeneratePrivKey(*algo, *bits)
		if err != nil {
			fmt.Printf("ERROR: Unable to generate key\n%v", err)
			os.Exit(1)
		}
	}

	if *genKey && !(*ephemeral) {
		if err = util.StorePrivKeyToFile(priv, *keyFile); err != nil {
			fmt.Printf("ERROR: Unable to save key to file %s\n", *keyFile)
			os.Exit(1)
		}
		fmt.Println("New key is stored at:", *keyFile)
	}

	if !(*ephemeral) {
		if !util.FileExists(*keyFile) {
			fmt.Printf("ERROR: Key (%s) does not exist.\n", *keyFile)
			fmt.Printf("Ensure path is correct or generate a new key with -genkey.\n")
			os.Exit(1)
		}

		priv, err = util.LoadPrivKeyFromFile(*keyFile)
		if err != nil {
			fmt.Printf("ERROR: Unable to load key from file\n%v", err)
			os.Exit(1)
		}
	}

	// TODO: Remove hard-coded listen address
	config := p2pnode.NewConfig()
	config.ListenAddrs = append(config.ListenAddrs, "/ip4/0.0.0.0/tcp/4001")
	config.PrivKey = priv
	config.BootstrapPeers = []string{} // We're bootsrapping, remove any default peers

	ctx := context.Background()
	_, err = p2pnode.NewNode(ctx, config)
	if err != nil {
		fmt.Println("ERROR: Unable to create new node\n", err)
		panic(err)
	}

	select {
	case <-ctx.Done(): // Likely will never happen...
		fmt.Println("ERROR: Main background context ended\n", ctx.Err())
		return
	}
}
