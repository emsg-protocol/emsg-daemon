# EMSG Daemon: Development Procedure Guide

This guide outlines the step-by-step procedure for implementing the EMSG Daemon, based on the project objectives and requirements.

---

## 1. DNS Routing Implementation (`router.go`) ✅ Complete
- Parse addresses in the format `user#domain.com`.
- Implement DNS TXT record lookup for `_emsg.domain.com`.
- Add unit tests for DNS routing logic.

## 2. Message Handling (`message.go`) ✅ Complete
- Implement logic for sending, receiving, and verifying messages.
- Add support for forwarding messages to group members using `group_id`.
- Integrate digital signature verification.
- Add unit tests for message validation and routing.

## 3. Persistent Storage (`storage.go`) ✅ Complete
- Implement CRUD operations for users, groups, and messages using SQLite/Postgres.
- Ensure storage supports both individual and group messages.
- Add unit tests for storage and retrieval.

## 4. Group Management & System Messages (`group.go`, `system.go`) ✅ Complete
- Expand group management to persist group state and membership.
- Implement system messages (group created, user joined/removed, etc.).
- Integrate system messages with group/user events.
- Add unit tests for group and system message logic.

## 5. Public Key Identity & Authentication (`auth.go`) ✅ Complete
- Integrate persistent storage for public keys.
- Add REST endpoints (optional) for key registration and management.
- Add unit tests for key registration and signature validation.

## 6. Configuration Loading (`config.go`) ✅ Complete
- Enhance config loading to support `.env` or config file in addition to environment variables.
- Add unit tests for config loading.

## 7. REST API (Optional, `api.go`) ✅ Complete
- Implement REST API endpoints for admin panel or client integration.
- Add unit tests for API endpoints.

## 8. Finalize & Expand Test Coverage
- Ensure all modules have comprehensive unit tests.
- Add integration tests for end-to-end flows (routing, messaging, group events).

---

**All core steps are now complete.**

Refer to this guide throughout the development process to ensure all requirements are met.
