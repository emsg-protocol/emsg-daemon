# 🛡️ AI Agent Prompt: Secure User Registration for EMSG Daemon

## 🎯 Objective

Enforce domain-level access control for user creation in the EMSG protocol. Prevent unauthorized users from registering under protected domains (e.g., `user#magadhaempire.com`) by introducing domain-based registration policies in the `emsg-daemon`.

---

## 🧩 Key Features to Implement

1. **Domain-specific user registration policies**:
   - `open`: Anyone can register (default)
   - `invite-only`: Requires valid invite code
   - `allowlist`: Only approved usernames can register
   - `admin-approval`: Registration request must be approved manually

2. **Invite code validation** for `invite-only` mode

3. **Allowlist lookup** for `allowlist` mode

4. **Pending queue storage** for `admin-approval` mode

5. **Domain policy loader**: Read from JSON/YAML file for each domain

6. **Enhanced API handler** in the daemon's registration endpoint to validate policy compliance

---

## 📁 Domain Policy Configuration Format

File: `config/domains/<domain>.json`

```json
{
  "domain": "magadhaempire.com",
  "mode": "invite-only",
  "invite_codes": ["MAJESTY2025", "TRUSTEDCIRCLE"],
  "allowlist": ["sandip", "admin"],
  "admin_notify": "admin@magadhaempire.com"
}
```

---

## 🔧 Backend Implementation Steps

- [ ] Add `DomainPolicy` struct and JSON loader
- [ ] Hook into registration flow (POST /register)
- [ ] Parse `user#domain.com` and fetch domain policy
- [ ] Validate request against policy
- [ ] On `admin-approval`, store request in `pending_registrations/`
- [ ] On failure, return clear error to client (e.g., `Invite code required`, `User not allowed`, `Pending approval`)

---

## 🔐 SDK Registration Input Format (Extended)

```json
{
  "address": "user#domain.com",
  "public_key": "base64...",
  "timestamp": 1234567890,
  "nonce": "random_string",
  "signature": "signed_payload",
  "invite_code": "MAJESTY2025"
}
```

> NOTE: `invite_code` is optional and depends on domain policy.

---

## ✉️ Optional Extensions

- Webhook for `admin_notify` on approval-needed request
- Signed invite codes (JWT or HMAC-based)
- Admin CLI for approving/rejecting pending requests

---

## 🧭 Files to Modify in `emsg-daemon`

- `handlers/register.go` (or equivalent endpoint handler)
- `config/policy_loader.go` (new or reused config reader)
- `models/policy.go` (new struct definitions)
- `pending/` (folder for pending approvals)
- `utils/validation.go` (error messages, input checks)

---

## 🔗 Related Projects

- [EMSG Daemon](https://github.com/emsg-protocol/emsg-daemon)
- [EMSG Client SDK](https://github.com/emsg-protocol/emsg-client-sdk)
- [EMSG Protocol Spec](https://github.com/emsg-protocol/specification)

This prompt is designed to guide the secure implementation of domain-level user registration with cryptographic validation and policy enforcement.