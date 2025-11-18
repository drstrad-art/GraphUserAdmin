# License Management Help

## Overview

This tool provides comprehensive license management for Microsoft 365 users and groups through the Microsoft Graph API.

## Available Commands

### List Available SKUs
```bash
gua licenses list-skus
```
Shows all available license SKUs in your tenant with their IDs and consumption.

### User License Management

#### View User Licenses
```bash
gua licenses get <UPN>
```
Example:
```bash
gua licenses get cbaker@alliance-hs.org
```

#### Add License to User
```bash
gua licenses add-user <UPN> <SKU_ID> [SKU_ID...]
```
Examples:
```bash
# Add single license
gua licenses add-user cbaker@alliance-hs.org c7df2760-2c81-4ef7-b578-5b5392b571df

# Add multiple licenses
gua licenses add-user cbaker@alliance-hs.org <SKU_ID_1> <SKU_ID_2>
```

#### Remove License from User
```bash
gua licenses remove-user <UPN> <SKU_ID> [SKU_ID...]
```
Examples:
```bash
# Remove single license
gua licenses remove-user cbaker@alliance-hs.org c7df2760-2c81-4ef7-b578-5b5392b571df

# Remove multiple licenses
gua licenses remove-user cbaker@alliance-hs.org <SKU_ID_1> <SKU_ID_2>
```

### Group License Management

Group-based licensing automatically assigns licenses to all members of a group.

#### View Group Licenses
```bash
gua licenses get-group <GROUP_ID>
```
Example:
```bash
gua licenses get-group a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

#### Add License to Group
```bash
gua licenses add-group <GROUP_ID> <SKU_ID> [SKU_ID...]
```
Examples:
```bash
# Add single license
gua licenses add-group a1b2c3d4-e5f6-7890-abcd-ef1234567890 c7df2760-2c81-4ef7-b578-5b5392b571df

# Add multiple licenses
gua licenses add-group <GROUP_ID> <SKU_ID_1> <SKU_ID_2>
```

#### Remove License from Group
```bash
gua licenses remove-group <GROUP_ID> <SKU_ID> [SKU_ID...]
```
Examples:
```bash
# Remove single license
gua licenses remove-group a1b2c3d4-e5f6-7890-abcd-ef1234567890 c7df2760-2c81-4ef7-b578-5b5392b571df

# Remove multiple licenses
gua licenses remove-group <GROUP_ID> <SKU_ID_1> <SKU_ID_2>
```

## Finding IDs

### How to Find SKU IDs

1. Run the list-skus command:
```bash
gua licenses list-skus
```

2. Output shows:
```
SKU Part Number          SKU ID                                  Consumed Units
ENTERPRISEPACK           c7df2760-2c81-4ef7-b578-5b5392b571df    25
POWER_BI_PRO            f8a1db68-be16-40ed-86d5-cb42ce701560    10
```

3. Use the SKU ID (the long UUID) in your commands

### Common SKU Part Numbers

| SKU Part Number | Description |
|----------------|-------------|
| ENTERPRISEPACK | Office 365 E3 |
| ENTERPRISEPREMIUM | Office 365 E5 |
| SPE_E3 | Microsoft 365 E3 |
| SPE_E5 | Microsoft 365 E5 |
| POWER_BI_PRO | Power BI Pro |
| FLOW_FREE | Power Automate Free |
| TEAMS_EXPLORATORY | Microsoft Teams Exploratory |
| PROJECTPROFESSIONAL | Project Plan 3 |
| VISIOCLIENT | Visio Plan 2 |

### How to Find Group IDs

**Method 1: Using gua tool**
```bash
gua groups get user@example.com
```

**Method 2: Azure Portal**
1. Go to Azure AD > Groups
2. Select the group
3. Copy the "Object ID"

## Complete Workflow Examples

### Example 1: Assign Office 365 E3 to a User
```bash
# Step 1: Find the SKU ID
gua licenses list-skus

# Step 2: Assign the license (using the SKU ID from step 1)
gua licenses add-user jdoe@example.com c7df2760-2c81-4ef7-b578-5b5392b571df

# Step 3: Verify
gua licenses get jdoe@example.com
```

### Example 2: Set Up Group-Based Licensing
```bash
# Step 1: Find your group ID
gua groups get admin@example.com

# Step 2: Assign license to the group
gua licenses add-group a1b2c3d4-e5f6-7890-abcd-ef1234567890 c7df2760-2c81-4ef7-b578-5b5392b571df

# Step 3: Verify
gua licenses get-group a1b2c3d4-e5f6-7890-abcd-ef1234567890

# All group members now have the license automatically!
```

### Example 3: Replace a User's License
```bash
# Check current licenses
gua licenses get user@example.com

# Remove old license (E3)
gua licenses remove-user user@example.com c7df2760-2c81-4ef7-b578-5b5392b571df

# Add new license (E5)
gua licenses add-user user@example.com 06ebc4ee-1bb5-47dd-8120-11324bc54e06

# Verify the change
gua licenses get user@example.com
```

### Example 4: Bulk License Assignment Using Groups
```bash
# Instead of assigning to individual users, use group-based licensing:

# 1. Get your group ID
gua groups get someuser@example.com

# 2. Assign licenses to the group
gua licenses add-group <GROUP_ID> <SKU_ID>

# Now all group members automatically have the license
# Adding users to the group automatically assigns the license
# Removing users automatically removes the license
```

## Troubleshooting

### Error: "Insufficient licenses available"
**Problem:** Not enough licenses purchased
**Solution:** 
- Run `gua licenses list-skus` to check available units
- Purchase more licenses in Microsoft 365 admin center

### Error: "User not found"
**Problem:** Invalid user principal name
**Solution:**
- Verify UPN with `gua users get <UPN>`
- Check for typos in email address

### Error: "Forbidden" or "Access denied"
**Problem:** Missing API permissions
**Solution:**
1. Go to Azure Portal > App Registrations
2. Select your app
3. Add these permissions:
   - `Directory.Read.All`
   - `Directory.ReadWrite.All`
   - `Organization.Read.All`
4. Grant admin consent

### Group-Based Licensing Not Working
**Problem:** License not assigned to group members
**Causes & Solutions:**
- **Usage location not set**: Users must have a usage location (e.g., "US", "GB")
- **Wrong group type**: Must be security group or Microsoft 365 group
- **License assignment errors**: Check Azure AD portal for specific errors

### Error: "No licenses assigned"
**Problem:** User/group has no licenses
**Solution:** This is informational, not an error. Assign licenses as needed.

## Best Practices

### 1. Use Group-Based Licensing
Instead of assigning licenses to individual users, create groups and assign licenses to groups:
- Easier management
- Automatic assignment/removal
- Better for large organizations

### 2. Always Check Before Removing
```bash
# Always verify current licenses first
gua licenses get user@example.com

# Then remove
gua licenses remove-user user@example.com <SKU_ID>
```

### 3. Keep a SKU Reference
Save your common SKU IDs for quick reference:
```bash
gua licenses list-skus > my-skus.txt
```

### 4. Monitor License Consumption
Regularly check usage:
```bash
gua licenses list-skus
```

### 5. Set Usage Location First
Before assigning licenses, ensure users have a usage location set (Microsoft requirement).

### 6. Test with One User
When making bulk changes, test with one user first:
```bash
# Test
gua licenses add-user test@example.com <SKU_ID>

# If successful, proceed with bulk operations
```

### 7. Use Configuration Files
For different environments:
```bash
gua --config prod-config.json licenses list-skus
gua --config test-config.json licenses list-skus
```

## Getting Help

### Show All Commands
```bash
gua --help
```

### Show License Commands
```bash
gua licenses --help
```

### Show Specific Command Help
```bash
gua licenses add-user --help
gua licenses add-group --help
```

## Quick Reference

| Task | Command |
|------|---------|
| List SKUs | `gua licenses list-skus` |
| Get user licenses | `gua licenses get <UPN>` |
| Add user license | `gua licenses add-user <UPN> <SKU_ID>` |
| Remove user license | `gua licenses remove-user <UPN> <SKU_ID>` |
| Get group licenses | `gua licenses get-group <GROUP_ID>` |
| Add group license | `gua licenses add-group <GROUP_ID> <SKU_ID>` |
| Remove group license | `gua licenses remove-group <GROUP_ID> <SKU_ID>` |
| Find groups | `gua groups get <UPN>` |
| Get user details | `gua users get <UPN>` |

## PowerShell Helper Script

For convenience, use the included helper script:

```powershell
# Show help
.\license.ps1 help

# List SKUs
.\license.ps1 list-skus

# User operations
.\license.ps1 user-get user@example.com
.\license.ps1 user-add user@example.com <SKU_ID>
.\license.ps1 user-remove user@example.com <SKU_ID>

# Group operations
.\license.ps1 group-get <GROUP_ID>
.\license.ps1 group-add <GROUP_ID> <SKU_ID>
.\license.ps1 group-remove <GROUP_ID> <SKU_ID>
```
