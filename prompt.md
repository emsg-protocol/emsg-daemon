
# 🧠 AI Agent Initialization Prompt for EMSG-Daemon

You are an advanced AI coding agent contributing to an open-source messaging server called the **EMSG Daemon** (Email Simplified Messaging Gateway).

Your task is to build a **decentralized messaging protocol** that replaces legacy email architecture with a modern, minimal, and secure communication protocol.

📌 Start by reading the file: `docs/EMSG_Copilot_AI_Agent_Guide.docx` in this repo.

## 🔹 Key Protocol Concepts:
- Users are addressed via `user#domain.com` instead of `user@domain.com`
- Messages consist of `to`, `cc`, `group_id`, and `body` fields only (no subject, bcc, MIME)
- A `group_id` enables persistent threads (like Telegram/Slack groups), across domains
- Public keys are used for authentication and digital signatures

## 🔹 Your Goals:
- Implement clean, modular Go code
- Support DNS TXT-based routing using `_emsg.domain.com`
- Build a persistent group system (`group.go`) for sending to all members
- Add a pubkey-based identity system in `auth.go` using Ed25519 or similar
- Use SQLite or Postgres for message and group storage
- Write reusable helper functions for routing, storage, and validation
- Add system messages like: group created, user joined, removed, etc.
- Optional: Build a REST API interface (`api.go`) for testing or client use

## 🔹 Suggested File Layout:
- `main.go`: Starts the daemon
- `router.go`: DNS-based routing
- `group.go`: Create & manage group state
- `auth.go`: Key registration & signature verification
- `message.go`: Core message handling
- `storage.go`: SQLite/Postgres handling
- `system.go`: Handles system-generated messages (e.g. "User joined")
- `config.go`: Environment and config loading

## 🔹 Testing Instructions:
- Write tests for DNS TXT lookups, pubkey verification, and group routing
- Ensure messages sent to a group are delivered to all members
- Validate error handling for missing or invalid `group_id`, `to`, or `body`

## 🔹 Example Message Format:
```json
{
  "from": "king#magadhaempire.com",
  "to": ["builder#emsg.dev"],
  "cc": ["group42#emsg.chat"],
  "group_id": "group42",
  "body": "Let's finalize the routing layer today.",
  "signature": "base64(sig)"
}
```

You must maintain high code quality, production readiness, and modularity in all implementations. Avoid SMTP/IMAP logic. This is not a forked mail server — it's a clean protocol layer.

Refer to `docs/EMSG_Copilot_AI_Agent_Guide.docx` for details.

➡️ Begin with scaffolding `main.go`, `router.go`, and `config.go`. Then implement `auth.go` and `group.go`.
