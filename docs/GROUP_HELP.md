# Group Management Help

## Overview

View Microsoft 365 group memberships through the Microsoft Graph API.

## Available Commands

### List All Groups
```bash
gua groups list
```
Shows all groups in your tenant with their display name, ID, and description.

Output example:
```
Display Name      ID                                      Description
------------      --                                      -----------
IT Department     a1b2c3d4-e5f6-7890-abcd-ef1234567890   IT Staff
Finance Team      b2c3d4e5-f6a7-8901-bcde-f12345678901   Finance Department
```

### Get User's Group Memberships
```bash
gua groups get <UPN>
```
Shows all groups that a user is a member of.

Example:
```bash
gua groups get cbaker@alliance-hs.org
```

Output example:
```
Display Name              ID
------------              --
Sales Team                a1b2c3d4-e5f6-7890-abcd-ef1234567890
Marketing Group           b2c3d4e5-f6a7-8901-bcde-fa1234567890
All Employees             c3d4e5f6-a7b8-9012-cdef-ab1234567890
```

### Add User to a Group
```bash
gua groups add-user <GROUP_ID> <UPN>
```
Example:
```bash
gua groups add-user a1b2c3d4-e5f6-7890-abcd-ef1234567890 jdoe@example.com
```

### Remove User from a Group
```bash
gua groups remove-user <GROUP_ID> <UPN>
```
Example:
```bash
gua groups remove-user a1b2c3d4-e5f6-7890-abcd-ef1234567890 jdoe@example.com
```

## Use Cases

### List All Groups
```bash
gua groups list
```

### Find Which Groups a User Belongs To
```bash
gua groups get user@example.com
```

### Get Group IDs for License Management
```bash
# Step 1: Find group IDs
gua groups get admin@example.com

# Step 2: Use group ID for license assignment
gua licenses add-group <GROUP_ID> <SKU_ID>
```

### Verify Group Membership
```bash
gua groups get user@example.com
```
Check if a user is in the expected groups.

## Finding Group IDs

### Method 1: Using GraphUserAdmin (gua)
```bash
gua groups get user@example.com
```
Copy the ID from the output.

### Method 2: Azure Portal
1. Go to Azure Active Directory
2. Select "Groups"
3. Find and select your group
4. Copy the "Object ID"

## Troubleshooting

### Error: "User not found"
- Verify the user principal name (UPN) is correct
- Run `gua users get <UPN>` to verify user exists

### Error: "Access denied"
- Your app registration needs these permissions:
  - `GroupMember.Read.All`
  - `Directory.Read.All`
- Grant admin consent in Azure Portal

### No Groups Shown
- User is not a member of any groups
- This is normal for new users
- Add user to groups in Azure Portal or PowerShell

## Quick Reference

| Task | Command |
|------|---------|
| List all groups | `gua groups list` |
| Get user's groups | `gua groups get <UPN>` |
| Add user to group | `gua groups add-user <GROUP_ID> <UPN>` |
| Remove user from group | `gua groups remove-user <GROUP_ID> <UPN>` |
| Get group ID for licenses | `gua groups get <UPN>` then copy ID |

## Related Commands

After getting a group ID, you can:

```bash
# View group licenses
gua licenses get-group <GROUP_ID>

# Add license to group
gua licenses add-group <GROUP_ID> <SKU_ID>

# Remove license from group
gua licenses remove-group <GROUP_ID> <SKU_ID>
```

See [LICENSE_HELP.md](LICENSE_HELP.md) for more details on license management.
