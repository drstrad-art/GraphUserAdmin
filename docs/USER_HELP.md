# User Management Help

## Overview

Manage Microsoft 365 users through the Microsoft Graph API.

## Available Commands

### List All Users
```bash
gua users list
```
Shows all users in your tenant with their display name, UPN, and email.

Output example:
```
Display Name      User Principal Name           Mail
------------      -------------------           ----
John Doe          jdoe@example.com             jdoe@example.com
Jane Smith        jsmith@example.com           jsmith@example.com
```

### Get User Details
```bash
gua users get <UPN>
```
Example:
```bash
gua users get cbaker@alliance-hs.org
```

Output shows:
- ID
- Display Name
- User Principal Name
- Mail
- Mail Nickname
- Account Enabled status

### Create a New User
```bash
gua users create <UPN> <DISPLAY_NAME> <MAIL_NICKNAME> <PASSWORD>
```
Example:
```bash
gua users create jdoe@example.com "John Doe" jdoe "TempPassword123!"
```

Notes:
- User will be required to change password on first sign-in
- Password must meet your tenant's complexity requirements

### Update User Properties
```bash
gua users update <UPN> <PROPERTY> <VALUE>
```
Examples:
```bash
gua users update jdoe@example.com displayName "Jane Doe"
gua users update jdoe@example.com usageLocation US
gua users update jdoe@example.com department "IT"
gua users update jdoe@example.com jobTitle "Senior Developer"
```

Common properties:
- displayName
- jobTitle
- department
- officeLocation
- mobilePhone
- usageLocation (required for license assignment)

### Delete a User
```bash
gua users delete <UPN>
```
Example:
```bash
gua users delete jdoe@example.com
```

Notes:
- You will be prompted for confirmation
- Deleted users can be restored within 30 days

## Examples

### Find a Specific User
```bash
gua users get jdoe@example.com
```

### List All Users and Save to File
```bash
gua users list > users.txt
```

### Create a New User with Required Properties
```bash
# Create the user
gua users create newuser@example.com "New User" newuser "InitialPass123!"

# Set usage location (required for license assignment)
gua users update newuser@example.com usageLocation US

# Assign a license
gua licenses add-user newuser@example.com <SKU_ID>
```

### Update Multiple Properties
```bash
gua users update jdoe@example.com displayName "Jane Doe"
gua users update jdoe@example.com department "Engineering"
gua users update jdoe@example.com jobTitle "Lead Developer"
```

## Troubleshooting

### Error: "User not found"
- Verify the user principal name (UPN) is correct
- Check for typos
- Ensure the user exists in your tenant

### Error: "Access denied"
- Your app registration needs these permissions:
  - `User.Read.All`
  - `Directory.Read.All`
- Grant admin consent in Azure Portal

## Quick Reference

| Task | Command |
|------|---------|
| List all users | `gua users list` |
| Get user details | `gua users get <UPN>` |
| Create user | `gua users create <UPN> <NAME> <NICKNAME> <PASSWORD>` |
| Update user | `gua users update <UPN> <PROPERTY> <VALUE>` |
| Delete user | `gua users delete <UPN>` |
