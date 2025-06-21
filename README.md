# EMSG Daemon

A decentralized messaging daemon implementing the EMSG (Electronic Message) protocol for secure, DNS-routed communication.

## Table of Contents

- [Overview](#overview)
- [EMSG Protocol Specification](#emsg-protocol-specification)
- [Architecture](#architecture)
- [Installation & Setup](#installation--setup)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Authentication](#authentication)
- [DNS Routing](#dns-routing)
- [Database Schema](#database-schema)
- [Development](#development)
- [Testing](#testing)
- [Production Deployment](#production-deployment)

## Overview

EMSG Daemon is a production-ready implementation of the EMSG protocol - a decentralized messaging system that uses DNS for service discovery and Ed25519 cryptographic signatures for authentication. Unlike traditional messaging systems that rely on centralized servers, EMSG enables direct server-to-server communication using DNS TXT records for routing.

### Key Features

- 🔐 **Ed25519 Cryptographic Authentication** - Secure message signing and verification
- 🌐 **DNS-Based Routing** - Decentralized service discovery via DNS TXT records
- 💾 **BoltDB Storage** - Pure Go database with no CGO dependencies
- 🔒 **Protected Endpoints** - Authentication middleware for sensitive operations
- 👥 **Group Messaging** - Multi-user group communication support
- ⚙️ **Environment Configuration** - Production-ready configuration system
- 🚀 **RESTful API** - Complete HTTP API for all operations

## EMSG Protocol Specification

### Address Format

EMSG addresses follow the format: `user#domain.com`

- **User Part**: Unique identifier within the domain
- **Domain Part**: DNS domain hosting the EMSG service
- **Separator**: `#` character (URL-encoded as `%23` in HTTP requests)

### DNS Service Discovery

EMSG uses DNS TXT records at `_emsg.domain.com` for service discovery:

```dns
_emsg.example.com. TXT "https://emsg.example.com:8080"
```

Or with structured JSON:

```dns
_emsg.example.com. TXT '{"server":"https://emsg.example.com:8080","version":"1.0","ttl":3600}'
```

### Message Format

Messages are JSON objects with the following structure:

```json
{
  "from": "alice#example.com",
  "to": ["bob#example.com"],
  "cc": ["charlie#example.com"],
  "group_id": "optional-group-id",
  "body": "Message content",
  "signature": "base64-ed25519-signature"
}
```

### Authentication Protocol

Authentication uses Ed25519 signatures with the following format:

```json
{
  "address": "user#domain.com",
  "timestamp": 1640995200,
  "nonce": "unique-value",
  "signature": "base64-ed25519-signature"
}
```

The signature covers: `METHOD:PATH:TIMESTAMP:NONCE`

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   EMSG Client   │    │   EMSG Daemon   │    │   DNS Server    │
│                 │    │                 │    │                 │
│ - Message UI    │◄──►│ - REST API      │◄──►│ - TXT Records   │
│ - Key Mgmt      │    │ - Auth Middleware│    │ - Service Disc. │
│ - Crypto        │    │ - Message Router│    │                 │
└─────────────────┘    │ - BoltDB Store  │    └─────────────────┘
                       │ - Group Mgmt    │
                       └─────────────────┘
                               │
                               ▼
                       ┌─────────────────┐
                       │   BoltDB        │
                       │                 │
                       │ - Users         │
                       │ - Messages      │
                       │ - Groups        │
                       └─────────────────┘
```

### Core Components

1. **REST API Server** - HTTP endpoints for all operations
2. **Authentication Middleware** - Ed25519 signature verification
3. **DNS Router** - Service discovery and message routing
4. **BoltDB Storage** - Persistent data storage
5. **Configuration System** - Environment-based configuration

### Technology Stack

- **Language**: Go 1.19+
- **Database**: BoltDB (pure Go, embedded)
- **Cryptography**: Ed25519 (crypto/ed25519)
- **HTTP Server**: Go standard library
- **DNS**: Go net package
- **Configuration**: Environment variables

## Installation & Setup

### Prerequisites

- Go 1.19 or later
- No CGO dependencies required

### Building from Source

```bash
# Clone the repository
git clone https://github.com/your-org/emsg-daemon.git
cd emsg-daemon

# Build the daemon
go build ./cmd/daemon

# Run tests
go test ./test/...
```

### Quick Start

```bash
# Run with default configuration
./daemon

# Run with custom configuration
export EMSG_DOMAIN="yourdomain.com"
export EMSG_PORT="8080"
export EMSG_DATABASE_URL="./emsg.db"
./daemon
```

The daemon will start on port 8080 (or your configured port) and create a BoltDB database file.

## Configuration

EMSG Daemon is configured via environment variables:

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `EMSG_DOMAIN` | `""` | Your domain name for EMSG service |
| `EMSG_DATABASE_URL` | `""` | Path to BoltDB database file |
| `EMSG_PORT` | `"8080"` | HTTP server port |
| `EMSG_LOG_LEVEL` | `"info"` | Logging level (debug, info, warn, error) |
| `EMSG_MAX_CONNECTIONS` | `100` | Maximum concurrent connections |

### Configuration Examples

**Development:**
```bash
export EMSG_DOMAIN="dev.emsg.local"
export EMSG_DATABASE_URL="./dev_emsg.db"
export EMSG_PORT="8080"
export EMSG_LOG_LEVEL="debug"
```

**Production:**
```bash
export EMSG_DOMAIN="emsg.yourdomain.com"
export EMSG_DATABASE_URL="/var/lib/emsg/emsg.db"
export EMSG_PORT="8080"
export EMSG_LOG_LEVEL="info"
export EMSG_MAX_CONNECTIONS="500"
```

**Windows PowerShell:**
```powershell
$env:EMSG_DOMAIN="emsg.yourdomain.com"
$env:EMSG_DATABASE_URL="C:\emsg\emsg.db"
$env:EMSG_PORT="8080"
```

## API Documentation

The EMSG Daemon provides a comprehensive REST API for all operations. Base URL: `http://localhost:8080` (or your configured port).

### User Management

#### Register User
```http
POST /api/user
Content-Type: application/json

{
  "address": "alice#example.com",
  "pubkey": "base64-encoded-ed25519-public-key",
  "first_name": "Alice",
  "middle_name": "",
  "last_name": "Smith",
  "display_picture": "https://example.com/alice.jpg"
}
```

**Response (201 Created):**
```json
{
  "address": "alice#example.com",
  "pubkey": "base64-encoded-ed25519-public-key",
  "first_name": "Alice",
  "last_name": "Smith",
  "display_picture": "https://example.com/alice.jpg"
}
```

#### Get User
```http
GET /api/user?address=alice%23example.com
```

**Response (200 OK):**
```json
{
  "address": "alice#example.com",
  "pubkey": "base64-encoded-ed25519-public-key",
  "first_name": "Alice",
  "last_name": "Smith",
  "display_picture": "https://example.com/alice.jpg"
}
```

### Message Management

#### Send Message (Protected)
```http
POST /api/message
Content-Type: application/json
Authorization: EMSG base64-encoded-auth-request

{
  "from": "alice#example.com",
  "to": ["bob#example.com"],
  "cc": ["charlie#example.com"],
  "group_id": "optional-group-id",
  "body": "Hello, this is a test message!",
  "signature": "base64-ed25519-signature"
}
```

**Response (201 Created):**
```json
{
  "status": "message sent"
}
```

#### Get Messages (Protected)
```http
GET /api/messages?user=alice%23example.com
Authorization: EMSG base64-encoded-auth-request
```

**Response (200 OK):**
```json
[
  {
    "from": "bob#example.com",
    "to": ["alice#example.com"],
    "cc": [],
    "group_id": "",
    "body": "Hello Alice!",
    "signature": "base64-ed25519-signature"
  }
]
```

### Group Management

#### Create Group (Protected)
```http
POST /api/group
Content-Type: application/json
Authorization: EMSG base64-encoded-auth-request

{
  "id": "dev-team",
  "name": "Development Team",
  "description": "EMSG Development Team Chat",
  "display_pic": "https://example.com/dev-team.jpg",
  "members": ["alice#example.com", "bob#example.com"]
}
```

**Response (201 Created):**
```json
{
  "ID": "dev-team",
  "Members": ["alice#example.com", "bob#example.com"],
  "Admins": null,
  "Name": "Development Team",
  "Description": "EMSG Development Team Chat",
  "DisplayPic": "https://example.com/dev-team.jpg"
}
```

#### Get Group
```http
GET /api/group?id=dev-team
```

**Response (200 OK):**
```json
{
  "ID": "dev-team",
  "Members": ["alice#example.com", "bob#example.com"],
  "Admins": null,
  "Name": "Development Team",
  "Description": "EMSG Development Team Chat",
  "DisplayPic": "https://example.com/dev-team.jpg"
}
```

### DNS Routing

#### Get Route Information
```http
GET /api/route?address=alice%23example.com
```

**Response (200 OK):**
```json
{
  "server": "https://emsg.example.com:8080",
  "pubkey": "base64-encoded-domain-public-key",
  "version": "1.0",
  "ttl": 3600
}
```

#### Validate Addresses
```http
POST /api/route/validate
Content-Type: application/json

{
  "addresses": [
    "alice#example.com",
    "invalid-address",
    "bob#test.com"
  ]
}
```

**Response (200 OK):**
```json
{
  "results": {
    "alice#example.com": {
      "valid": true
    },
    "invalid-address": {
      "valid": false,
      "error": "invalid address format: must be user#domain.com"
    },
    "bob#test.com": {
      "valid": true
    }
  }
}
```

#### Route Message
```http
POST /api/route/message
Content-Type: application/json

{
  "recipients": [
    "alice#example.com",
    "bob#test.com"
  ]
}
```

**Response (200 OK):**
```json
{
  "routes": {
    "https://emsg.example.com:8080": ["alice#example.com"],
    "https://emsg.test.com:8080": ["bob#test.com"]
  }
}
```

### Error Responses

All endpoints return appropriate HTTP status codes and error messages:

**400 Bad Request:**
```json
{
  "error": "invalid request format"
}
```

**401 Unauthorized:**
```json
{
  "error": "authentication failed: signature verification failed"
}
```

**404 Not Found:**
```json
{
  "error": "user not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "internal server error"
}
```

## Authentication

EMSG uses Ed25519 cryptographic signatures for authentication. Protected endpoints require an `Authorization` header.

### Authorization Header Format

```
Authorization: EMSG <base64-encoded-auth-request>
```

### Auth Request Structure

The auth request is a base64-encoded JSON object:

```json
{
  "address": "alice#example.com",
  "timestamp": 1640995200,
  "nonce": "unique-random-value",
  "signature": "base64-ed25519-signature"
}
```

### Signature Generation

The signature is generated over the string: `METHOD:PATH:TIMESTAMP:NONCE`

Example for `POST /api/message`:
```
POST:/api/message:1640995200:abc123
```

### Implementation Example (Go)

```go
import (
    "crypto/ed25519"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "time"
)

func createAuthRequest(address string, privateKey ed25519.PrivateKey, method, path string) (string, error) {
    timestamp := time.Now().Unix()
    nonce := fmt.Sprintf("%d", time.Now().UnixNano())

    // Create message to sign
    message := fmt.Sprintf("%s:%s:%d:%s", method, path, timestamp, nonce)

    // Sign the message
    signature := ed25519.Sign(privateKey, []byte(message))

    // Create auth request
    authReq := map[string]interface{}{
        "address":   address,
        "timestamp": timestamp,
        "nonce":     nonce,
        "signature": base64.StdEncoding.EncodeToString(signature),
    }

    // Encode as JSON then base64
    jsonData, _ := json.Marshal(authReq)
    return base64.StdEncoding.EncodeToString(jsonData), nil
}
```

### Protected Endpoints

The following endpoints require authentication:

- `POST /api/message` - Send messages
- `GET /api/messages` - Retrieve messages
- `POST /api/group` - Create groups

### Security Features

- **Timestamp Validation**: Prevents replay attacks (5 min past, 1 min future window)
- **Nonce**: Prevents duplicate requests
- **Ed25519 Signatures**: Cryptographically secure authentication
- **Public Key Verification**: User's public key must be registered

## DNS Routing

EMSG uses DNS TXT records for decentralized service discovery and message routing.

### DNS Record Setup

To enable EMSG for your domain, create a TXT record at `_emsg.yourdomain.com`:

**Simple Format:**
```dns
_emsg.example.com. 3600 IN TXT "https://emsg.example.com:8080"
```

**Structured Format:**
```dns
_emsg.example.com. 3600 IN TXT '{"server":"https://emsg.example.com:8080","version":"1.0","ttl":3600,"pubkey":"domain-public-key"}'
```

### Routing Process

1. **Address Parsing**: Extract domain from `user#domain.com`
2. **DNS Lookup**: Query `_emsg.domain.com` TXT record
3. **Route Parsing**: Parse server URL from TXT record
4. **Message Delivery**: Send message to discovered server

### Example DNS Configuration

For domain `example.com` running EMSG on `emsg.example.com:8080`:

```dns
; EMSG service discovery
_emsg.example.com. 3600 IN TXT "https://emsg.example.com:8080"

; Optional: Domain public key for verification
_emsg.example.com. 3600 IN TXT "pubkey=base64-encoded-ed25519-public-key"
```

### Routing Implementation

The daemon automatically handles routing for outbound messages:

```go
// Example: Sending to alice#example.com
// 1. Lookup _emsg.example.com TXT record
// 2. Parse server URL: https://emsg.example.com:8080
// 3. POST message to https://emsg.example.com:8080/api/message
```

## Database Schema

EMSG Daemon uses BoltDB with the following bucket structure:

### Users Bucket
```
Key: "alice#example.com"
Value: {
  "address": "alice#example.com",
  "pubkey": "base64-ed25519-public-key",
  "first_name": "Alice",
  "middle_name": "",
  "last_name": "Smith",
  "display_picture": "https://example.com/alice.jpg"
}
```

### Messages Bucket
```
Key: "alice#example.com" (recipient)
Value: [
  {
    "from": "bob#example.com",
    "to": ["alice#example.com"],
    "cc": [],
    "group_id": "",
    "body": "Hello Alice!",
    "signature": "base64-ed25519-signature"
  }
]
```

### Groups Bucket
```
Key: "dev-team"
Value: {
  "ID": "dev-team",
  "Members": ["alice#example.com", "bob#example.com"],
  "Admins": ["alice#example.com"],
  "Name": "Development Team",
  "Description": "EMSG Development Team Chat",
  "DisplayPic": "https://example.com/dev-team.jpg"
}
```

### Database Operations

- **Users**: Create, Read (no Update/Delete for security)
- **Messages**: Create, Read (append-only for audit trail)
- **Groups**: Create, Read, Update (membership changes)

## Development

### Project Structure

```
emsg-daemon/
├── cmd/daemon/          # Main application entry point
├── api/                 # REST API handlers and middleware
├── internal/
│   ├── auth/           # Authentication utilities
│   ├── config/         # Configuration management
│   ├── group/          # Group management
│   ├── message/        # Message handling
│   ├── router/         # DNS routing logic
│   ├── storage/        # Database operations
│   └── system/         # System utilities
├── test/               # Test files
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # This file
```

### Key Packages

- **`cmd/daemon`**: Main application with server startup
- **`api`**: HTTP handlers and authentication middleware
- **`internal/storage`**: BoltDB operations and data models
- **`internal/router`**: DNS TXT record lookup and routing
- **`internal/auth`**: Ed25519 signature verification
- **`internal/config`**: Environment variable configuration

### Building and Running

```bash
# Install dependencies
go mod download

# Build the daemon
go build -o daemon ./cmd/daemon

# Run with default settings
./daemon

# Run with custom configuration
EMSG_PORT=9090 EMSG_DOMAIN=test.local ./daemon
```

## Testing

The project includes comprehensive test suites for all components.

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test files
go run test_user_api.go
go run test_message_group_api.go
go run test_dns_routing.go
go run test_auth_middleware.go
go run test_env_config.go

# Final integration test
go run test_final_integration.go
```

### Test Coverage

- **User Management**: Registration, retrieval, validation
- **Message System**: Sending, receiving, authentication
- **Group Management**: Creation, membership, retrieval
- **DNS Routing**: Address validation, route lookup, message routing
- **Authentication**: Signature verification, timestamp validation
- **Configuration**: Environment variables, defaults

### Test Files

- `test_user_api.go` - User registration and retrieval
- `test_message_group_api.go` - Message and group operations
- `test_dns_routing.go` - DNS routing functionality
- `test_auth_middleware.go` - Authentication middleware
- `test_env_config.go` - Environment configuration
- `test_final_integration.go` - Complete integration test

### Manual Testing

```bash
# Start daemon
./daemon

# Test user registration
curl -X POST http://localhost:8080/api/user \
  -H "Content-Type: application/json" \
  -d '{"address":"test#example.com","pubkey":"...","first_name":"Test"}'

# Test user retrieval
curl "http://localhost:8080/api/user?address=test%23example.com"

# Test address validation
curl -X POST http://localhost:8080/api/route/validate \
  -H "Content-Type: application/json" \
  -d '{"addresses":["test#example.com","invalid"]}'
```

## Production Deployment

### System Requirements

- **OS**: Linux, macOS, Windows
- **Memory**: 512MB RAM minimum, 2GB recommended
- **Storage**: 1GB minimum for database and logs
- **Network**: Port 8080 (or configured port) accessible
- **DNS**: TXT record at `_emsg.yourdomain.com`

### Deployment Steps

1. **Build the Binary**
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o emsg-daemon ./cmd/daemon
   ```

2. **Set Environment Variables**
   ```bash
   export EMSG_DOMAIN="emsg.yourdomain.com"
   export EMSG_DATABASE_URL="/var/lib/emsg/emsg.db"
   export EMSG_PORT="8080"
   export EMSG_LOG_LEVEL="info"
   export EMSG_MAX_CONNECTIONS="1000"
   ```

3. **Create Database Directory**
   ```bash
   sudo mkdir -p /var/lib/emsg
   sudo chown emsg:emsg /var/lib/emsg
   ```

4. **Configure DNS**
   ```dns
   _emsg.yourdomain.com. 3600 IN TXT "https://emsg.yourdomain.com:8080"
   ```

5. **Start the Service**
   ```bash
   ./emsg-daemon
   ```

### Systemd Service

Create `/etc/systemd/system/emsg-daemon.service`:

```ini
[Unit]
Description=EMSG Daemon
After=network.target

[Service]
Type=simple
User=emsg
Group=emsg
WorkingDirectory=/opt/emsg
ExecStart=/opt/emsg/emsg-daemon
Environment=EMSG_DOMAIN=emsg.yourdomain.com
Environment=EMSG_DATABASE_URL=/var/lib/emsg/emsg.db
Environment=EMSG_PORT=8080
Environment=EMSG_LOG_LEVEL=info
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable emsg-daemon
sudo systemctl start emsg-daemon
sudo systemctl status emsg-daemon
```

### Docker Deployment

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o emsg-daemon ./cmd/daemon

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/emsg-daemon .
EXPOSE 8080
CMD ["./emsg-daemon"]
```

Build and run:
```bash
docker build -t emsg-daemon .
docker run -d \
  -p 8080:8080 \
  -e EMSG_DOMAIN=emsg.yourdomain.com \
  -e EMSG_DATABASE_URL=/data/emsg.db \
  -v /var/lib/emsg:/data \
  emsg-daemon
```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name emsg.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSL/TLS Configuration

For production, use HTTPS with Let's Encrypt:

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d emsg.yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Monitoring and Logging

- **Logs**: Monitor daemon output for errors
- **Health Check**: `GET /api/user?address=health` (should return 404)
- **Metrics**: Monitor response times and error rates
- **Database**: Monitor BoltDB file size and performance

### Security Considerations

- **Firewall**: Only expose necessary ports (80, 443, 8080)
- **User Permissions**: Run daemon as non-root user
- **Database Security**: Protect database file permissions
- **DNS Security**: Use DNSSEC for TXT records
- **Rate Limiting**: Implement rate limiting for API endpoints

### Backup and Recovery

```bash
# Backup database
cp /var/lib/emsg/emsg.db /backup/emsg-$(date +%Y%m%d).db

# Restore database
systemctl stop emsg-daemon
cp /backup/emsg-20231201.db /var/lib/emsg/emsg.db
systemctl start emsg-daemon
```

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Support

For questions and support:
- Create an issue on GitHub
- Check the documentation
- Review the test files for examples

---

**EMSG Daemon** - Decentralized messaging for the modern web.
