# Distributed File Storage Using Golang 
A Scalable, Decentralized, Fully distributed, Content addressable file storage system using Golang that can handle and stream very large files.

## Project Structure
```
.
├── Makefile
├── README.md
├── crypto.go             - This file defines methods to encrypt files using AES encryption.
├── crypto_test.go        - This file contains unit tests for the crypto.go file.
├── go.mod
├── go.sum
├── learn.exe             - This is a binary executable file that can be run directly without the Golang compiler.
├── main.go               - This file contains the configuration for file, crypto, and p2p.
├── p2p
│   ├── encoding.go             - This file specifies the encoding to efficiently transfer files.
│   ├── handshake.go            - This file contains pre-connection checks. Additional tests can be added before establishing a connection.
│   ├── message.go              - This file defines the message type and format, such as whether it is a message or a stream of files.
│   ├── tcp_transport.go        - This file defines methods to handle messages between two nodes.
│   ├── tcp_transport_test.go   - This file contains unit tests for the tcp_transport.go file.
│   └── transport.go            - This file defines structs and interfaces for peer and transport, respectively.
├── server.go             - This file defines methods to store and retrieve files from clients and nodes.
├── storage.go            - This file defines methods to store, retrieve, and check files on the local disk.
└── storage_test.go       - This file contains unit tests for the storage.go file.

2 directories, 18 files
```

> [!NOTE]
>- This repository is ideal for those who want to learn Golang at its peak. It offers a great opportunity to explore and understand how such systems work under the hood.
>- While more features could be added, such as automatic peer discovery, automatic file recovery, and load balancers for faster file retrieval, this project is already a perfect starting point for beginners. Keep in mind that the goal is to learn and explore, and there will always be more features to add.
>- New features and pull requests are welcome.

### Features:
- Concurrency: Utilize goroutines to efficiently handle multiple operations.
- Scalability: Designed for horizontal scaling.
- Fault Tolerance: Implements data replication and failover mechanisms.
- Security: Ensures data encryption and secure communications.
- Performance: Optimized for low latency and high throughput.

### High-Level Diagram of Distributed File Storage Using Golang:
![image](https://github.com/user-attachments/assets/7405a81a-bbed-44cd-a09e-0e3d443ba87f)

Feel free to explore and contribute to this project!
