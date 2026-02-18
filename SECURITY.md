# Security Policy

## Supported Versions

Security fixes are provided for the latest version on the `main` branch only.

| Version        | Supported |
| -------------- | --------- |
| Latest         | Yes       |
| Older versions | No        |

---

## Reporting a Vulnerability

If you discover a security vulnerability, please report it privately using **GitHub Security Advisories**:

Repository → Security → Advisories → “Report a vulnerability”

Do not create public issues for security vulnerabilities.

Please provide:

* A detailed description of the issue
* Steps to reproduce or proof-of-concept
* Impact assessment
* Any suggested remediation (optional)

All reports will be reviewed and addressed as quickly as possible.

---

## Scope

We welcome reports related to:

### Authentication and Accounts

* Authentication bypass
* Password reset vulnerabilities
* Session fixation or hijacking
* Privilege escalation

### User Data and Privacy

* Unauthorized access to user profiles or applications
* Insecure Direct Object References (IDOR)
* Exposure of sensitive or personal data

### API and Backend Security

* Injection vulnerabilities (SQL, NoSQL, command)
* Mass assignment issues
* Improper access control
* Insufficient rate limiting

### Content and File Handling

* Stored or reflected XSS in job postings
* Malicious file upload vulnerabilities
* HTML injection

### Secrets and Infrastructure

* Hardcoded credentials or tokens
* Misconfigured storage or services
* Sensitive data exposure

---

## Out of Scope

The following are typically not considered valid reports:

* Missing best-practice headers without exploitability
* Self-XSS or clickjacking on public pages
* Issues requiring physical access
* Theoretical attacks without proof-of-concept
* Vulnerabilities in third-party services without a demonstrated impact on this project

---

## Responsible Disclosure

Please:

* Avoid accessing real user data
* Avoid service disruption or denial-of-service testing
* Keep findings confidential until a fix is released

---
