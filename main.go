// // package main

// // import (
// // 	"context"
// // 	"fmt"
// // 	"os"
// // 	"os/signal"
// // 	"syscall"

// // 	"github.com/libp2p/go-libp2p"
// // 	peerstore "github.com/libp2p/go-libp2p-peerstore"
// // 	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
// // )

// // func main() {
// // 	// create a background context (i.e. one that never cancels)
// // 	ctx := context.Background()

// // 	// start a libp2p node with default settings
// // 	node, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
// // 		libp2p.Ping(false),)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	// configure our own ping protocol
// // 	pingService := &ping.PingService{Host: node}
// // 	node.SetStreamHandler(ping.ID, pingService.PingHandler)

// // 	// print the node's listening addresses
// // 	// fmt.Println("Listen addresses:", node.Addrs())
// // 	// print the node's PeerInfo in multiaddr format
// // 	peerInfo := &peerstore.PeerInfo{
// // 		ID:    node.ID(),
// // 		Addrs: node.Addrs(),
// // 	}
// // 	addrs, err := peerstore.InfoToP2pAddrs(peerInfo)
// // 	if err != nil {
// // 		panic(err)
// // 	}
// // 	fmt.Println("libp2p node address:", addrs[0])

// //     // wait for a SIGINT or SIGTERM signal
// //     ch := make(chan os.Signal, 1)
// //     signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
// //     <-ch
// //     fmt.Println("Received signal, shutting down...")


// // 	// shut the node down
// // 	if err := node.Close(); err != nil {
// // 		panic(err)
// // 	}
// // }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/libp2p/go-libp2p"
// 	peerstore "github.com/libp2p/go-libp2p-peerstore"
// 	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
// 	multiaddr "github.com/multiformats/go-multiaddr"
// )

// func main() {
// 	// create a background context (i.e. one that never cancels)
// 	ctx := context.Background()
// 	var DefaultPeerstore Option = func(cfg *Config) error {
// 	    return cfg.Apply(Peerstore(pstoremem.NewPeerstore("QmZJexJchxXt71N9bSkj6PchxyGEs9g6Qvrj7fEk32FqDs")))
// 	}

// 	// peerID := ID ("QmSgHqe9daPQ6fUKgAjJNPsD5fQNEMwTYmcEJ3ptUyF3kE")
// 	// peerID := ID ("QmZJexJchxXt71N9bSkj6PchxyGEs9g6Qvrj7fEk32FqDs")

// 	// start a libp2p node that listens on a random local TCP port,
// 	// but without running the built-in ping protocol
// 	node, err := libp2p.New(ctx, libp2p.Identity(DefaultPeerstore), 
// 		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
// 		libp2p.Ping(false),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// configure our own ping protocol
// 	pingService := &ping.PingService{Host: node}
// 	node.SetStreamHandler(ping.ID, pingService.PingHandler)

// 	// print the node's PeerInfo in multiaddr format
// 	peerInfo := &peerstore.PeerInfo{
// 		ID:    node.ID(),
// 		Addrs: node.Addrs(),
// 	}
// 	fmt.Println("PEER INFO:", peerInfo)
// 	addrs, err := peerstore.InfoToP2pAddrs(peerInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("libp2p node address:", addrs[0])

// 	// if a remote peer has been passed on the command line, connect to it
// 	// and send it 5 ping messages, otherwise wait for a signal to stop
// 	if len(os.Args) > 1 {
// 		addr, err := multiaddr.NewMultiaddr(os.Args[1])
// 		if err != nil {
// 			panic(err)
// 		}
// 		peer, err := peerstore.InfoFromP2pAddr(addr)
// 		if err != nil {
// 			panic(err)
// 		}
// 		if err := node.Connect(ctx, *peer); err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("sending 5 ping messages to", addr)
// 		ch := pingService.Ping(ctx, peer.ID)
// 		for i := 0; i < 5; i++ {
// 			res := <-ch
// 			fmt.Println("pinged", addr, "in", res.RTT)
// 		}
// 	} else {
// 		// wait for a SIGINT or SIGTERM signal
// 		ch := make(chan os.Signal, 1)
// 		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
// 		<-ch
// 		fmt.Println("Received signal, shutting down...")
// 	}

// 	// shut the node down
// 	if err := node.Close(); err != nil {
// 		panic(err)
// 	}
// }



package main

import (
	"context"
	"fmt"
	"flag"
	// "os"
	// "os/signal"
	// "syscall"
	// Option "option"
	// "rsa"

	"github.com/libp2p/go-libp2p-core/crypto"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	// "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	// create a background context (i.e. one that never cancels)
	// type Config = config.Config
	// type Option = config.Option
	ctx := context.Background()
	// help := flag.Bool("help", false, "Display Help")
	listenHost := flag.String("host", "127.0.0.1", "The bootstrap node host listen address\n")
	port := flag.Int("port", 4001, "The bootstrap node listen port")
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

	// start a libp2p node with default settings

	// priv := "&{[211 151 127 224 159 14 157 18 23 132 211 171 4 8 125 131 235 83 169 205 79 230 32 138 150 179 103 28 152 240 11 111 101 134 246 174 231 186 183 172 59 180 89 156 126 43 240 153 190 62 31 24 209 96 245 188 19 240 39 95 93 41 140 38]}"
	byteVal := []byte {211, 151, 127, 224, 159, 14, 157, 18, 23, 132, 211, 171, 4, 8, 125, 131, 235, 83, 169, 205, 79, 230, 32, 138, 150, 179, 103, 28, 152, 240, 11, 111, 101, 134, 246, 174, 231, 186, 183, 172, 59, 180, 89, 156, 126, 43, 240, 153, 190, 62, 31, 24, 209, 96, 245, 188, 19, 240, 39, 95, 93, 41, 140, 38}
	priv, _ := crypto.UnmarshalEd25519PrivateKey(byteVal)
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", *listenHost, *port))
	node, err := libp2p.New(ctx, libp2p.ListenAddrs(sourceMultiAddr), libp2p.Identity(priv), )
	if err != nil {
		panic(err)
	}
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("This node: ", node.ID().Pretty(), " ", node.Addrs())
	_, err = dht.New(ctx, node)

	peerInfo := &peerstore.PeerInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	fmt.Println("PEER INFO:", peerInfo)

	// print the node's listening addresses
	// fmt.Println("Listen addresses:", node.Addrs())



	// shut the node down
	// if err := node.Close(); err != nil {
	// 	panic(err)
	// }
	select{}
}