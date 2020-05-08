# Bootstrap Node
Launches a P2P node to serve as a bootstrap node for others P2P nodes to connect to

Supports different cryptographic algorithms for creating the private/public key pair.

Run `bootstrap -h` to see the the full list of CLI flags and options:
```
Usage of ./bootstrap:
  -algo string
        Cryptographic algorithm to use for generating the key.
        Will be ignored if 'genkey' is false.
        Must be one of {RSA, Ed25519, Secp256k1, ECDSA} (default "RSA")
  -bits int
        Key length, in bits. Will be ignored if 'algo' is not RSA. (default 2048)
  -ephemeral
        Generate a new key just for this run, and don't store it to file.
        If 'keyfile' is specified, it will be ignored.
  -genkey
        Generate a new key and save to file.
  -keyfile string
        Location of private key to read from (or write to, if generating). (default "~/.privKey")
```

