package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func main() {
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


	byteVal := []byte {211, 151, 127, 224, 159, 14, 157, 18, 23, 132, 211, 171, 4, 8, 125, 131, 235, 83, 169, 205, 79, 230, 32, 138, 150, 179, 103, 28, 152, 240, 11, 111, 101, 134, 246, 174, 231, 186, 183, 172, 59, 180, 89, 156, 126, 43, 240, 153, 190, 62, 31, 24, 209, 96, 245, 188, 19, 240, 39, 95, 93, 41, 140, 38}
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

	select{}
}
