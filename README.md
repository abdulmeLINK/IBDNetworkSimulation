# IBDNetworkSimulation

## Description

This project is a simulation of the Bitcoin Core Initial Block Download (IBD) process. It consists of a client and multiple servers. The client requests blocks from the servers, and the servers respond with block data.

## Getting Started

### Prerequisites

- Go 1.16 or later

### Installation

1. Clone the repository:

```sh
git clone https://github.com/abdulmeLINK/IBDNetworkSimulation
```

2. Navigate to the project directory:

```sh
cd project
```

3. Install the dependencies:

```sh
go mod download
```

## Usage

### Running the Servers

To start the servers, navigate to the `cmd/server` directory and run:

```sh
go run main.go <port> <block_info.txt file location>
```

### Running the Client

To start the client, navigate to the `cmd/client` directory and run:

```sh
go run main.go "idrangestart-idrangeend"
```

## Testing

To run the tests, navigate to the project root directory and run:

```sh
scripts/start_simulation.sh
```

