# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| 0.4.x   | Yes       |
| 0.3.x   | Security fixes only |
| < 0.3   | No        |

## Reporting a vulnerability

Please report vulnerabilities privately to **security@finflow.example**. Include a description, reproduction steps, and the affected component (`apps/web`, `services/ledger`, `services/fraud`, or `apps/worker`). Encrypt sensitive details with our PGP key if possible.

We aim to:

- acknowledge your report within 2 business days,
- provide an initial assessment within 5 business days,
- ship a fix or mitigation for high-severity issues within 30 days.

Please do not open public issues or discuss the vulnerability publicly until a fix has been released.

## Scope

In scope: authentication and session handling, authorization checks on accounts and transfers, ledger integrity, idempotency handling, webhook signing, and fraud-rule bypasses.

Out of scope: denial of service through volumetric traffic, social engineering of maintainers, and findings that require a compromised end-user device.

## Handling secrets

Never commit credentials. Use `.env` locally (it is git-ignored) and Kubernetes secrets in deployed environments. If a secret is exposed, rotate it immediately and notify the security contact.
