# 🔄 AI Agent Prompt: Admin-Based User Identity Reset for EMSG Daemon

## 🎯 Objective

Implement a secure admin-driven recovery/reset system in the `emsg-daemon`, allowing domain administrators to reset or update a user's identity (e.g., when a private key is lost or compromised).

---

## 🔐 Use Case

Only **authorized domain administrators** (defined by public keys in the domain config) can:

- Reset a user's public key
- Deactivate or delete a user (optional extension)

---

## 📄 Endpoint

**POST /admin/reset_user**

### Request Fields

```json
{
  "operation": "reset_user",
  "address": "sandip#magadhaempire.com",
  "new_public_key": "base64-encoded-new-pubkey",
  "timestamp": 1723812738,
  "nonce": "random_nonce",
  "admin_signature": "base64-signature"
}
```

---

## ✅ Validation Logic

1. Parse `address` → `user#domain.com`
2. Load domain policy from `config/domains/<domain>.json`
3. Check if request is signed by one of `admin_public_keys`
4. Validate `timestamp` (must be recent ±5min)
5. Verify `nonce` hasn't been reused (store/cache)
6. Replace the existing user’s public key with `new_public_key`
7. Log the operation and return a success response

---

## 🛡️ Domain Policy Example (JSON)

```json
{
  "domain": "magadhaempire.com",
  "mode": "admin-approval",
  "admin_public_keys": [
    "base64-admin-pubkey-1",
    "base64-admin-pubkey-2"
  ]
}
```

---

## 🔏 Signature Rules

Signature = `Sign(admin_private_key, address + new_public_key + timestamp + nonce)`

Use Ed25519 with `crypto/ed25519` Go standard library.

---

## 🧾 Response Format

| Status Code | Description                      |
|-------------|----------------------------------|
| 200 OK      | Public key successfully updated  |
| 403 Forbidden | Invalid signature or not admin |
| 409 Conflict | User does not exist             |
| 429 Too Many Requests | Rate-limited           |

---

## 🧰 Optional Enhancements

- Maintain history of previous public keys (for audit)
- Notify previous user key via system message or webhook
- Add CLI utility `./emsgctl reset-user ...`
- Write reset actions to `audit.log`

---

## 📦 Files To Modify

- `handlers/admin.go` (create this if not exists)
- `models/domain_policy.go`
- `utils/crypto.go` (signature verify, nonce check)
- `config/domains/<domain>.json`
- `audit/audit.go` (optional)

---

## 🔗 Related Projects

- [EMSG Daemon](https://github.com/emsg-protocol/emsg-daemon)
- [EMSG Client SDK](https://github.com/emsg-protocol/emsg-client-sdk)

This feature allows safe user key recovery, matching enterprise-grade messaging standards.