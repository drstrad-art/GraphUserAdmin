# GraphUserAdmin (gua)

A command-line tool for managing Microsoft 365 users, licenses, and groups using the Microsoft Graph API with client credentials authentication.

## Features

- **User Management**: List, create, update, and delete Microsoft 365 users
- **License Management**: Assign, remove, and manage licenses for users and groups
- **Group Management**: View group memberships and manage group assignments

## Building

```powershell
.\build.ps1
```

Creates `./bin/gua.exe`

## Setup

1. Copy `config.json.example` to `config.json`
2. Fill in your Azure credentials:
   ```json
   {
     "tenantId": "your-tenant-id",
     "clientId": "your-client-id",
     "clientSecret": "your-client-secret"
   }
   ```

## Quick Start

```bash
# Show all commands
gua --help

# List users
gua users list

# List available licenses
gua licenses list-skus

# Add license to user
gua licenses add-user user@example.com <SKU_ID>

# List groups
gua groups list

# Add user to group
gua groups add-user <GROUP_ID> user@example.com
```

## Usage

```bash
gua [FLAGS] COMMAND [SUBCOMMAND] [OPTIONS]
```

**Global Flags:**
- `--config, -c` - Path to config file (default: config.json)
- `--verbose, -v` - Enable verbose output
- `--version` - Show version

**Commands:**
- `users` - Manage users (list, get, create, update, delete)
- `licenses` - Manage licenses (list-skus, get, add-user, remove-user, add-group, remove-group)
- `groups` - Manage groups (list, get, add-user, remove-user)

## Requirements

- Go 1.21+
- Azure app registration with appropriate Microsoft Graph API permissions
- Tenant ID, Client ID, and Client Secret

## Documentation

See `docs/` for detailed guides:
- `INDEX.md` - Complete command reference
- `USER_HELP.md` - User management guide
- `LICENSE_HELP.md` - License management guide
- `GROUP_HELP.md` - Group management guide

## License

See LICENSE file for details
