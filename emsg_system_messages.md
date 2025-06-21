# EMSG Protocol System Messages

## Overview
The EMSG protocol uses various system messages for communication, error handling, status reporting, and automated notifications.

## 📧 Message Types

### 1. User Messages (Application Level)
```json
{
  "from": "alice#example.com",
  "to": ["bob#example.com"],
  "cc": ["charlie#example.com"],
  "group_id": "optional-group-id",
  "body": "User message content",
  "signature": "base64-ed25519-signature"
}
```

### 2. System Messages (Protocol Level)
```json
{
  "from": "system#local",
  "to": ["group-id"],
  "cc": [],
  "group_id": "group-id",
  "body": "[SYSTEM] event_type: user alice#example.com in group dev-team",
  "signature": ""
}
```

## 🔧 System Message Types

### Group Management Events
```
SystemGroupCreated       = "group_created"
SystemUserJoined         = "user_joined"
SystemUserLeft           = "user_left"
SystemUserRemoved        = "user_removed"
SystemAdminAssigned      = "admin_assigned"
SystemAdminRevoked       = "admin_revoked"
SystemGroupRenamed       = "group_renamed"
SystemDescriptionUpdated = "description_updated"
SystemDPUpdated          = "dp_updated"
```

### Example System Messages
```
[SYSTEM] group_created: user alice#example.com in group dev-team
[SYSTEM] user_joined: user bob#example.com in group dev-team
[SYSTEM] admin_assigned: user alice#example.com in group dev-team
[SYSTEM] group_renamed: user alice#example.com in group dev-team
```

## 🌐 API Response Messages

### Success Responses
```json
// Message sent successfully
{
  "status": "message sent"
}

// User registered successfully
{
  "address": "alice#example.com",
  "pubkey": "base64-ed25519-public-key",
  "first_name": "Alice",
  "last_name": "Smith"
}

// Address validation results
{
  "results": {
    "alice#example.com": {
      "valid": true
    },
    "invalid-address": {
      "valid": false,
      "error": "invalid address format: must be user#domain.com"
    }
  }
}

// Message routing results
{
  "routes": {
    "https://emsg.example.com:8765": ["alice#example.com"],
    "https://emsg.test.com:8765": ["bob#test.com"]
  }
}
```

### Error Messages

#### 400 Bad Request
```
"missing address parameter"
"missing user parameter"
"invalid request"
"no recipients provided"
"missing required fields: from, to, or body"
"invalid address format: must be user#domain.com"
```

#### 401 Unauthorized
```
"missing authorization header"
"invalid authorization format"
"invalid authorization encoding"
"authentication failed: signature verification failed"
"authentication failed: timestamp out of range"
"authentication failed: user not found"
"authentication failed: invalid signature encoding"
```

#### 404 Not Found
```
"user not found"
"group not found"
"route lookup failed: no such host"
```

#### 500 Internal Server Error
```
"internal server error"
"database error"
"failed to store message"
"failed to create group"
```

## 🔐 Authentication Messages

### Auth Request Format
```json
{
  "address": "alice#example.com",
  "timestamp": 1640995200,
  "nonce": "unique-random-value",
  "signature": "base64-ed25519-signature"
}
```

### Auth Error Messages
```
"timestamp out of range"
"user not found: <error details>"
"invalid signature encoding"
"signature verification failed"
```

## 🌍 DNS and Routing Messages

### DNS TXT Record Formats
```
// Simple format
"https://emsg.example.com:8765"

// Structured format
{
  "server": "https://emsg.example.com:8765",
  "version": "1.0",
  "ttl": 3600,
  "pubkey": "optional-domain-public-key"
}
```

### Routing Error Messages
```
"route lookup failed: no such host"
"routing failed: invalid recipient address"
"unable to parse route info: <record>"
"no server found for <address>"
"failed to get route info for <address>"
```

## 📱 Client-Server Communication

### Request Headers
```
Content-Type: application/json
Authorization: EMSG <base64-encoded-auth-request>
```

### Response Headers
```
Content-Type: application/json
X-EMSG-User: <authenticated-user-address>
```

## 🔍 Validation Messages

### Address Validation
```
"invalid address format: must be user#domain.com"
"user part cannot be empty"
"domain part cannot be empty"
"domain must contain at least one dot"
```

### Message Validation
```
"missing required fields: from, to, or body"
"invalid recipient address format"
"signature verification failed"
```

## 📊 Status and Logging Messages

### Server Status
```
"Starting REST API server on :8765..."
"EMSG Daemon started successfully"
"Database initialized"
"Configuration loaded"
```

### Debug Messages
```
"DNS TXT lookup for: _emsg.example.com"
"Route found: https://emsg.example.com:8765"
"User authenticated: alice#example.com"
"Message stored for: bob#example.com"
```

## 🚨 Error Handling

### Network Errors
```
"connection timeout"
"DNS resolution failed"
"server unreachable"
"SSL certificate error"
```

### Database Errors
```
"database connection failed"
"user already exists"
"group already exists"
"failed to store message"
"failed to retrieve messages"
```

## 🔧 Configuration Messages

### Environment Variables
```
EMSG_DOMAIN="yourdomain.com"
EMSG_DATABASE_URL="./emsg.db"
EMSG_PORT="8765"
EMSG_LOG_LEVEL="info"
EMSG_MAX_CONNECTIONS="1000"
```

### Config Validation
```
"invalid port number"
"database file not found"
"invalid log level"
"domain name required"
```

## 📋 Message Flow Examples

### Successful Message Flow
1. Client: `POST /api/message` with auth
2. Server: Validates auth → `200 OK`
3. Server: Validates message → `200 OK`
4. Server: Stores message → `201 Created`
5. Response: `{"status": "message sent"}`

### Failed Authentication Flow
1. Client: `POST /api/message` without auth
2. Server: `401 Unauthorized`
3. Response: `"missing authorization header"`

### Cross-Domain Routing Flow
1. Client: Send to `bob#external.com`
2. Server: DNS lookup `_emsg.external.com`
3. Server: Route to `https://emsg.external.com:8765`
4. Response: Message delivered or routing error

## 🎯 Protocol Constants

### Message Types
```
USER_MESSAGE = "user_message"
SYSTEM_MESSAGE = "system_message"
GROUP_MESSAGE = "group_message"
```

### Event Types
```
MESSAGE_SENT = "message_sent"
MESSAGE_RECEIVED = "message_received"
USER_ONLINE = "user_online"
USER_OFFLINE = "user_offline"
```

This comprehensive list covers all system messages used in the EMSG protocol implementation.
