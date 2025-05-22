# VXET-SDK

VXET-SDK is the official development toolkit for the VXET blockchain, providing a complete set of tools for interacting with the VXET blockchain.

## Features

- Complete VXET blockchain client implementation
- Various command-line tools for development and debugging
- Rich API interfaces for application development
- Flexible consensus mechanism support
- Comprehensive RPC services

## System Requirements

- Go 1.20 or higher
- Compatible with Linux, MacOS, and Windows systems

## Installation Guide

### Building from Source

1. Clone the repository

```bash
git clone https://github.com/VXETChain/VXET-SDK.git
cd VXET-SDK
```

2. Compile the VXET client

```bash
make vxet
```

After compilation, the executable will be located in the `build/bin/` directory.

### Using Pre-compiled Versions

Visit the [Releases page](https://github.com/VXETChain/VXET-SDK/releases) to download the pre-compiled version suitable for your system.

## Usage

### Starting the VXET Client

```bash
./build/bin/vxet
```

### Development Tools

Install the required development tools:

```bash
make devtools
```

### Running Tests

```bash
make test
```

## Project Structure

- `cmd/`: Contains various command-line tools
  - `vxet/`: VXET client
  - `abigen/`: Smart contract ABI generation tool
  - Other utility tools
- `common/`: Common utilities and data structures
  - `hexutil/`: Hexadecimal utilities
  - `math/`: Mathematical operation tools
  - `prque/`: Priority queue implementation
- `consensus/`: Consensus algorithm implementation
  - `clique/`: PoA consensus algorithm
- `vxetclient/`: Client implementation
  - `gethclient/`: Ethereum client compatibility interface
- `vxet/`: Core blockchain implementation
  - `protocols/`: Network protocol implementation
  - `gasprice/`: Gas price management
  - `downloader/`: Block downloader
- `params/`: Network parameter configuration
- `rpc/`: RPC service implementation
- `p2p/`: P2P network layer implementation
- `vxetdb/`: Blockchain database
- `node/`: Node services
- `console/`: Console interface
- `graphql/`: GraphQL API implementation
- `signer/`: Transaction signing tools

## License

This project is licensed under GPL-3.0, please refer to the [COPYING](COPYING) file for details.

## Contribution Guidelines

Pull Requests and Issues are welcome. Before submitting code, please ensure that all tests are passed and follow the project's code standards.

## Contact Information

- Official website: [VXET Chain](https://www.vertexect.com/)
- Feedback: Please submit an Issue on GitHub 