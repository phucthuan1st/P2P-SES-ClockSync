# P2P Network with Clock Synchronization using Go

## Introduction

This repository contains an implementation of a P2P (Peer-to-Peer) network using the Go programming language. The project also includes the SES (Schiper-Eggli-Sandoz) Algorithm for clock synchronization among the peers in the network.

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [P2P Network](#p2p-network)
- [Clock Synchronization](#clock-synchronization)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

To get started with this project, follow these steps:

1. **Clone the Repository**: You can clone this repository to your local machine using Git.

```bash
git clone https://github.com/phucthuan1st/p2p-clock-sync.git
```

2. **Install Go**: Make sure you have Go installed on your system. You can download it from the [official Go website](https://golang.org/dl/).

3. **Build and Run**: Build and run the P2P network and clock synchronization system.

```bash
cd p2p-clock-sync
go run main.go --config config/NodeX.json
```

## P2P Network

This project simulates a P2P network, where nodes communicate with each other directly without the need for a central server. The network allows peer nodes to exchange data and synchronize their clocks using the SES Algorithm.

- **Node Communication**: Nodes in the P2P network communicate with each other using Go's networking features.

- **Data Exchange**: Nodes can exchange data, share information, and collaborate in a decentralized manner.

## Clock Synchronization

The project implements the SES (Source-Initiated Time Synchronization) Algorithm for clock synchronization among the P2P network nodes. This algorithm allows nodes to synchronize their clocks with a reference node.

- **SES Algorithm**: The Schiper-Eggli-Sandoz Algorithm is used to establish a common time reference among the network nodes.

- **Clock Precision**: The algorithm takes into account clock precision, network delays, and clock drift to ensure accurate time synchronization.

- **Collaborative Synchronization**: Nodes work together to achieve a synchronized network time, allowing for coordinated actions and timestamp accuracy.

## Contributing

We welcome contributions from the open-source community. If you have suggestions, improvements, or bug fixes, please feel free to create issues or submit pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
