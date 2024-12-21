<p align="center">
<img src="https://github.com/andygeiss/cloud-native-store/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Store

[![License](https://img.shields.io/github/license/andygeiss/cloud-native-store)](https://github.com/andygeiss/cloud-native-store/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-store)](https://github.com/andygeiss/cloud-native-store/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-store)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-store)

**Cloud Native Store** is a Go-based key-value store showcasing cloud-native patterns like transactional logging, data sharding, encryption, TLS, and circuit breakers. Built with hexagonal architecture for modularity and extensibility, it includes a robust API and in-memory storage for efficiency and stability.

## Project Motivation

The motivation behind **Cloud Native Store** is to provide a practical example of implementing a key-value store that adheres to cloud-native principles. The project aims to:

1. Demonstrate best practices for building scalable, secure, and reliable cloud-native applications.
2. Showcase the use of hexagonal architecture to enable modular and testable code.
3. Offer a reference implementation for features like encryption, transactional logging, and stability mechanisms.
4. Inspire developers to adopt cloud-native patterns in their projects.

## Project Setup and Run Instructions

Follow these steps to set up and run the **Cloud Native Store**:

### Prerequisites
1. Install [Go](https://go.dev/) (version 1.18 or higher).
2. Install [mkcert](https://github.com/FiloSottile/mkcert) for generating local TLS certificates.
3. Create an `.env` file with the following contents:

```env
ENCRYPTION_KEY="0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
SERVER_CERTIFICATE=".tls/server.crt"
SERVER_KEY=".tls/server.key"
TRANSACTIONAL_LOG=".cache/transactions.json"
```

### Commands

#### Run the Service
To start the service:
```bash
just run
```

#### Set Up the Service
To set up the necessary environment, including generating local TLS certificates:
```bash
just setup
```
This command will:
- Install `mkcert` (via Homebrew).
- Create the directories `.cache` and `.tls`.
- Generate a self-signed certificate for `localhost`.

#### Test the Service
To run unit tests:
```bash
just test
```
This will execute tests for the core service logic.

---

### Directory Structure Overview

```
.
├── LICENSE               # License file
├── README.md             # Documentation file
├── cmd/
│   └── main.go          # Entry point of the application
├── go.mod                # Go module definition
├── go.sum                # Go dependencies lock file
├── internal/
│   └── app/
│       ├── adapters/    # Adapters for inbound and outbound communication
│       │   ├── inbound/
│       │   │   └── api/ # HTTP API handlers and router
│       │   │       ├── handlers.go
│       │   │       └── router.go
│       │   └── outbound/
│       │       └── inmemory/ # In-memory object storage implementation
│       │           └── object_store.go
│       ├── config/      # Configuration handling
│       │   └── config.go
│       └── core/        # Core domain logic
│           ├── ports/   # Interfaces for core abstractions
│           │   └── object.go
│           └── services/ # Business logic implementations
│               ├── mocks_test.go
│               ├── object.go
│               └── object_test.go
└── logo.png              # Project logo
```

---

## How `main.go` Puts Everything Together

The `main.go` file is the entry point of the application, orchestrating the various components of the **Cloud Native Store** to create a functional and secure server. Below is a breakdown of its key aspects:

### Configuration Setup
The application reads and sets up necessary configurations:
- Encryption key for data security.
- Paths to TLS certificate and key files for secure communication.
- Path for the transactional log file.

Code snippet:
```go
cfg := &config.Config{
    Key:    security.Getenv("ENCRYPTION_KEY"),
    Server: config.Server{CertFile: os.Getenv("SERVER_CERTIFICATE"),
    KeyFile: os.Getenv("SERVER_KEY")},
}
```

### Transactional Logger Initialization
A transactional logger is created to ensure all operations are logged for auditability and consistency. The logger uses a JSON file for persistent event storage and processes events asynchronously.

Key features:
- **Buffered Channels**: Events are queued in a buffered channel for efficient processing.
- **Event Sequence Management**: Ensures all events are uniquely identified and ordered.
- **Graceful Shutdown**: Pending events are processed before shutting down the logger.

Code snippet:
```go
logger := consistency.NewJsonFileLogger[string, string](os.Getenv("TRANSACTIONAL_LOG"))
```

Behind the scenes:
- The logger reads the last sequence number from the log file during initialization.
- A dedicated goroutine processes events from the queue and writes them to the log file.
- Errors are captured and reported through an error channel.

### Object Storage and Service Setup
The in-memory object store is initialized and integrated with the core object service, which is also configured with the transactional logger:

Key features of Object Storage:
- **Sharding**: The in-memory object store employs sharding to optimize access and reduce contention.
- **CRUD Operations**: Supports create, read, update, and delete operations with thread safety.
- **Error Handling**: Returns specific errors, such as `ErrorKeyDoesNotExist`, for missing keys.

Key features of Object Service:
- **Stability Patterns**: Implements timeout, debounce, retry, and circuit breaker patterns for robustness.
- **Transactional Logging**: Logs all operations (put/delete) for auditability and recovery.
- **Encryption/Decryption**: Secures data during transmission and storage.

Code snippet:
```go
port := inmemory.NewObjectStore(1)
service := services.
    NewObjectService(cfg).
    WithTransactionalLogger(logger).
    WithPort(port)

if err := service.Setup(); err != nil {
    log.Fatalf("error during setup: %v", err)
}
```

Behind the scenes:
- **Initialization**: The object store initializes shards for efficient data partitioning.
- **Setup Processing**: The object service processes pending events from the transactional logger to ensure data consistency.
- **Lifecycle Management**: Teardown ensures graceful cleanup of resources.

### API Routing
The API router is initialized, binding HTTP endpoints to the core service. Each route is associated with a specific HTTP method and operation:

- **PUT `/api/v1/store`**: Adds or updates an object. Expects a JSON request body with `key` and `value` fields. Returns a status code of 200 on success.
- **GET `/api/v1/store`**: Retrieves an object by key. Expects a JSON request body with the `key` field. Returns the associated `value` or a 404 status if the key is not found.
- **DELETE `/api/v1/store`**: Deletes an object by key. Expects a JSON request body with the `key` field. Returns a status code of 200 on success.

Code snippet:
```go
mux := api.Route(service)
```

Behind the scenes:
- **JSON Parsing**: Handlers decode incoming requests and encode responses in JSON format.
- **Error Handling**: Handlers return appropriate HTTP status codes for invalid requests or errors.
- **Integration**: Handlers interact with the core object service to execute CRUD operations.

### Secure Server Initialization
A secure server is set up with advanced TLS settings to ensure encrypted communication and robust security:

Key features:
- **TLS Configuration**: Automatically acquires and manages TLS certificates for specified domains using Let's Encrypt. Supports a self-signed certificate for localhost.
- **Secure Cipher Suites**: Uses strong, modern ciphers for secure communication.
- **ALPN Support**: Enables application layer protocol negotiation (e.g., HTTP/2).
- **Timeouts**: Configures read, write, and idle timeouts to enhance server responsiveness and mitigate resource exhaustion attacks.

Code snippet:
```go
srv := security.NewServer(mux, "localhost")
log.Printf("start listening...")
if err := srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
    log.Fatalf("listening failed: %v", err)
}
```

Behind the scenes:
- **Dynamic Certificate Management**: The server uses `autocert.Manager` for automatic certificate acquisition and renewal.
- **Security Enhancements**: Implements secure defaults, including minimum TLS version 1.2, preferred server cipher suites, and limited max header size.
- **Graceful Handling**: Ensures proper setup and teardown of resources, preventing abrupt failures.

### Resource Cleanup
Proper teardown mechanisms are used to release resources and close connections:
```go
defer service.Teardown()
defer srv.Close()
```

By following this flow, `main.go` ties together the cloud-native patterns and components into a cohesive and functional key-value store application.
