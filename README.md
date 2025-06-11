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

