# GraphUserAdmin (gua) - Microsoft 365 Management Tool

## Quick Start

```bash
# List available commands
gua --help

# Check version
gua --version

# List users
gua users list

# Get user details
gua users get user@example.com

# Create a new user
gua users create user@example.com "User Name" username "TempPass123!"

# Update user property
gua users update user@example.com usageLocation US

# List available licenses
gua licenses list-skus

# Get user's licenses
gua licenses get user@example.com

# Add license to user
gua licenses add-user user@example.com <SKU_ID>

# Remove license from user
gua licenses remove-user user@example.com <SKU_ID>

# List all groups
gua groups list

# Get user's groups
gua groups get user@example.com

# Add user to group
gua groups add-user <GROUP_ID> user@example.com

# Add license to group
gua licenses add-group <GROUP_ID> <SKU_ID>
```

## Documentation

### Detailed Help by Topic

- **[License Management](LICENSE_HELP.md)** - Complete guide for managing user and group licenses
- **[User Management](USER_HELP.md)** - Guide for viewing and managing users
- **[Group Management](GROUP_HELP.md)** - Guide for viewing group memberships

### Quick Reference Guides

- **[LICENSE_MANAGEMENT_GUIDE.md](../LICENSE_MANAGEMENT_GUIDE.md)** - Quick reference for license operations
- **[README.md](../README.md)** - Main project documentation with setup instructions

## Common Tasks

### View Available Licenses in Your Tenant
```bash
gua licenses list-skus
```

### Assign a License to a User
```bash
# Step 1: Find the SKU ID
gua licenses list-skus

# Step 2: Assign it
gua licenses add-user user@example.com <SKU_ID>

# Step 3: Verify
gua licenses get user@example.com
```

### Set Up Group-Based Licensing
```bash
# Step 1: Find the group ID
gua groups get admin@example.com

# Step 2: Assign license to group
gua licenses add-group <GROUP_ID> <SKU_ID>

# Step 3: Verify
gua licenses get-group <GROUP_ID>
```

### Remove a License from a User
```bash
# Step 1: Check what licenses they have
gua licenses get user@example.com

# Step 2: Remove specific license
gua licenses remove-user user@example.com <SKU_ID>
```

## Getting More Help

### Show All Commands
```bash
gua --help
```

### Show Help for Specific Command Group
```bash
gua users --help
gua licenses --help
gua groups --help
```

### Show Help for Specific Command
```bash
gua licenses add-user --help
gua licenses add-group --help
```

## Helper Scripts

### PowerShell License Helper
```powershell
# Show help
.\license.ps1 help

# Quick operations
.\license.ps1 list-skus
.\license.ps1 user-add user@example.com <SKU_ID>
.\license.ps1 group-add <GROUP_ID> <SKU_ID>
```

## Configuration

### Default Config File
The tool looks for `config.json` in the current directory by default.

### Use Different Config File
```bash
gua --config /path/to/config.json users list
```

### Config File Format
```json
{
  "tenantId": "your-tenant-id",
  "clientId": "your-client-id",
  "clientSecret": "your-client-secret"
}
```

## Support & Troubleshooting

### Common Issues

**Authentication Errors**
- Verify credentials in config.json
- Check API permissions in Azure Portal
- Ensure admin consent is granted

**Permission Errors**
- See [README.md](../README.md) for required permissions
- Grant admin consent in Azure Portal

**License Assignment Errors**
- Check available license units with `gua licenses list-skus`
- Verify user has usage location set
- See [LICENSE_HELP.md](LICENSE_HELP.md) for detailed troubleshooting

## Quick Command Reference

| Category | Command | Description |
|----------|---------|-------------|
| **Users** | `gua users list` | List all users |
|  | `gua users get <UPN>` | Get user details |
|  | `gua users create <UPN> <NAME> <NICKNAME> <PASS>` | Create new user |
|  | `gua users update <UPN> <PROP> <VALUE>` | Update user property |
|  | `gua users delete <UPN>` | Delete user |
| **Licenses - View** | `gua licenses list-skus` | List available SKUs |
|  | `gua licenses get <UPN>` | Get user's licenses |
|  | `gua licenses get-group <ID>` | Get group's licenses |
| **Licenses - User** | `gua licenses add-user <UPN> <SKU>` | Add license to user |
|  | `gua licenses remove-user <UPN> <SKU>` | Remove license from user |
| **Licenses - Group** | `gua licenses add-group <ID> <SKU>` | Add license to group |
|  | `gua licenses remove-group <ID> <SKU>` | Remove license from group |
| **Groups** | `gua groups list` | List all groups |
|  | `gua groups get <UPN>` | Get user's groups |
|  | `gua groups add-user <ID> <UPN>` | Add user to group |
|  | `gua groups remove-user <ID> <UPN>` | Remove user from group |
| **General** | `gua --help` | Show all commands |
|  | `gua --version` | Show version |
|  | `gua --verbose <command>` | Enable debug output |
|  | `gua <command> --help` | Show command help |

## Next Steps

1. Read the [LICENSE_HELP.md](LICENSE_HELP.md) for comprehensive license management guide
2. Check [README.md](../README.md) for setup and configuration
3. Use `gua --help` to explore all available commands
