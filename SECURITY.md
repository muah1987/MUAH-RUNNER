# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| `main` branch | ✅ |
| Released tags | ✅ latest only |
| Older releases | ❌ |

---

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Email: create a private security advisory at **GitHub → Security → Advisories → New draft advisory**
(or email the maintainer directly if that is not available for this repository).

Include:
- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fix (optional but welcome)

You will receive an acknowledgment within **48 hours** and a fix timeline within **7 days** for critical issues.

---

## Security Architecture

### Chain-of-Thought (CoT) Encryption

All chain-of-thought reasoning produced by muah-runner is **encrypted at rest** using AES-256-GCM before being written to disk:

- **Key:** auto-generated on `muah-runner init`, stored in `.muah/config/secrets.enc`
- **Vault:** `.muah/privacy_cot/vault/` — every entry is individually encrypted
- **Classification:** each CoT entry is classified as:
  - `private` — never shown, not redactable
  - `redactable` — shown only with PII/secrets stripped
  - `public` — always available

```bash
muah-runner cot status         # Check encryption status
muah-runner cot redacted <id>  # View redacted CoT for a run
```

### Secrets Management

- Secrets are **never stored in plaintext** — always in `.muah/config/secrets.enc`
- API keys/tokens are read from environment variables only
- Never commit `.muah/` to version control (it is in `.gitignore` by default)
- The `muah-runner.yml` config file should never contain actual secret values — use env var references

### Network Security

- The REST API (`localhost:8080` by default) binds to `0.0.0.0` — restrict this in production
- GitHub webhook signatures are verified with HMAC-SHA256 (`webhooks.secret` config key)
- All outbound HTTPS calls use the system CA bundle

### Container Security

- The Docker image uses `alpine:3.21` (pinned major version) — not `latest`
- The binary is built with `-s -w` to strip debug symbols and reduce attack surface
- CGO is disabled (`CGO_ENABLED=0`) for a fully static binary
- The container does not run as root in production deployments — consider adding `USER nonroot`

### Dependency Management

- Go modules with `go.sum` hash verification
- Dependencies are pinned in `go.mod`
- Run `go mod tidy && go mod verify` before submitting PRs

---

## Known Limitations

- TLS is not yet implemented for the REST API (`server.tls: false` in config) — use a reverse proxy (nginx, Caddy) in front of muah-runner in production
- The Docker executor mounts the host Docker socket — this is a privileged operation; only use in trusted environments

---

## Security Hardening Checklist (for production deployments)

- [ ] Bind API to `127.0.0.1` or use a reverse proxy
- [ ] Set `webhooks.secret` to a strong random value
- [ ] Mount `.muah/` on encrypted storage
- [ ] Rotate the CoT encryption key periodically (`muah-runner cot status`)
- [ ] Use Docker secrets or a vault for API tokens (not `.env` files)
- [ ] Pin the Docker image to a specific digest in production
- [ ] Run the container as a non-root user
