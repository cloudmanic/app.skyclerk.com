# Skyclerk.com

## Development Requirements

### Frontend
- Node.js version 12.22.2 (required for Angular 7.2 compatibility)
  - Use `nvm use 12.22.2` if you have nvm installed
  - Or download directly from [nodejs.org](https://nodejs.org/)

## Backend Commands

### Purge Old Accounts

The `purge-old-accounts` command removes inactive accounts to keep the database clean. This is useful for cleaning up free trial accounts that were never used.

```bash
go run main.go --purge-old-accounts
```

#### Deletion Criteria

The command will delete accounts that meet ANY of these criteria:
- Accounts with zero ledger entries that were created over 6 months ago
- Accounts that haven't had any ledger activity in over one year

#### Protection Rules

Accounts are protected from deletion if:
- Any user associated with the account has the email address `spicer@cloudmanic.com`

#### What Gets Deleted

When an account is deleted, the following data is removed:
- All ledger entries and associated file/label relationships
- All activities, invites, labels, files, contacts, categories
- SnapClerk entries and connected accounts
- The account record itself
- Associated billing records (if not shared with other accounts)
- User records (if not associated with other accounts)
- Stripe customer data (if applicable)

#### Output

The command provides a detailed summary showing:
- Total accounts checked
- Accounts with zero ledgers
- Accounts with old ledgers
- Total accounts deleted
- Accounts skipped due to email protection rule

## Production deployment and registration security

GitHub Actions (`.github/workflows/fly-deploy.yml`) tests and deploys `master` to
Fly.io app `app-skyclerk-com`. The Digital Ocean instructions below are historical.

Set GitHub repository secrets `TURNSTILE_SITE_KEY` and `TURNSTILE_SECRET_KEY` from
1Password item **Cloudflare Turnstile - Skyclerk.com**. The existing `FLY_API_TOKEN`
authorizes deployment. The workflow requires both Turnstile keys and stages them
as Fly.io runtime secrets before deploying. No credentials are committed or
embedded in the frontend build. `/registration-config` exposes only the public key.

Turnstile protects both new accounts and invited registrations. Cloudflare must
return a successful verification with action `register` and hostname matching
`SITE_DOMAIN` (`app.skyclerk.com` in production). Ensure the widget allows
`app.skyclerk.com` in Cloudflare before deployment. Missing configuration or a
verification outage blocks signup. Failed submissions request a fresh challenge.

For end-to-end local development, use a separate development widget that allows
`localhost`, put its keys in `backend/.env`, and set `SITE_DOMAIN=localhost`.
Cloudflare's [dummy keys](https://developers.cloudflare.com/turnstile/troubleshooting/testing/)
can test the widget UI, but their fixed `test` action does not pass this backend's
strict `register` action check. There is no `APP_ENV` bypass. Automated tests mock
verification and never use real Turnstile keys.

The public website's support form uses the same widget and runtime secret.
`POST /support/contact-us` requires `turnstile_token`, verifies action
`support_contact`, and only accepts hostnames `skyclerk.com` and
`www.skyclerk.com`. The widget must allow the marketing domain in addition to
`app.skyclerk.com`. Requests with missing, rejected, or expired challenges cannot
send email or Slack notifications. Verification failures return HTTP 400 with
`code: "turnstile_failed"`; email delivery failures return HTTP 503 so the form
can preserve the message and request a fresh challenge for retry.

Deploy the companion marketing-site changes from `../skyclerk.com` before this
backend update so the public form sends tokens when verification becomes mandatory.
No additional secrets are needed beyond the existing Turnstile configuration.
Support tests live in `support_test.go`, with shared HTTP mocks in
`turnstile_test.go`; the existing registration tests also exercise the shared verifier.

# Deploying Servers

* When deploying a server with Digital Ocean copy the following into the `User-Data` filed. It will run Cloud Init when the VPS boots up.

```
#cloud-config
users:
  - name: spicer
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh-authorized-keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAsw21gLc2CaJN8hJB7zWDYWLF5gqWl6t8ozgso8aOrq8rz7P8ji3MwvHEelEe6UMNg4CxWTGYIWvFptlfCRvy9d94RBy9AAdb4pEBmSOyxPf8sJ+xD+V3TFJfmMOAm4049cBLN9b7+PRkUjl4jC3zTch5tQ+5lG7v04tWwzCaSCSD2HNuw2qKK3FpaLA6EIw+ieueBkgNgRnwMvgVO8nmyOkR5b3WUoL4vow3heNHV00V4M0yhBHLHDIFkXMgMztpLm3Dki1ZplUF0EyPH5llj5a4n2RMR5c7B1wAiXuUPO0oQTw9ItS5SZl9zKu9ZuIvqeXWsz/0NqRdEMIKqvxIZQ== spicer@cloudmanic.com
packages:
  - python2
```

* Once a fresh server is up and running configure it with `ansible-playbook server-config.yml`
