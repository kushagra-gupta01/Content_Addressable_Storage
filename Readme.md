# Content Addressable Storage

This is a custom distributed CAS memory implementation in golang.

For establishing connection and facilitating communication between peers a custom library  **p2p library** has been implemented.

A file once sent n the network will be replicated to all peers.