package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"GraphUserAdmin/internal/groups"
	"GraphUserAdmin/internal/licenses"
	"GraphUserAdmin/internal/users"

	"github.com/spf13/cobra"
)

// setupCommands creates and configures all CLI commands
func setupCommands(rootCmd *cobra.Command) {
	setupUsersCommands(rootCmd)
	setupLicensesCommands(rootCmd)
	setupGroupsCommands(rootCmd)
}

// setupUsersCommands creates the users command and its subcommands
func setupUsersCommands(rootCmd *cobra.Command) {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Manage Microsoft 365 users",
	}

	usersListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all users in the tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			userList, err := users.ListUsers(token, "")
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "Display Name\tUser Principal Name\tMail")
			fmt.Fprintln(w, "------------\t-------------------\t----")
			for _, user := range userList {
				fmt.Fprintf(w, "%s\t%s\t%s\n", user.DisplayName, user.UserPrincipalName, user.Mail)
			}
			w.Flush()

			return nil
		},
	}

	usersGetCmd := &cobra.Command{
		Use:   "get [UPN]",
		Short: "Get details for a specific user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := users.GetUser(token, args[0])
			if err != nil {
				return err
			}

			fmt.Printf("ID:                  %s\n", user.ID)
			fmt.Printf("Display Name:        %s\n", user.DisplayName)
			fmt.Printf("User Principal Name: %s\n", user.UserPrincipalName)
			fmt.Printf("Mail:                %s\n", user.Mail)
			fmt.Printf("Mail Nickname:       %s\n", user.MailNickname)
			fmt.Printf("Account Enabled:     %t\n", user.AccountEnabled)

			return nil
		},
	}

	usersCreateCmd := &cobra.Command{
		Use:   "create [UPN] [DISPLAY_NAME] [MAIL_NICKNAME] [PASSWORD]",
		Short: "Create a new user",
		Long:  "Create a new user with the specified details. The user will be required to change password on first sign-in.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			upn := args[0]
			displayName := args[1]
			mailNickname := args[2]
			password := args[3]

			user, err := users.CreateUser(token, displayName, upn, mailNickname, password, true)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully created user!\n")
			fmt.Printf("ID:                  %s\n", user.ID)
			fmt.Printf("Display Name:        %s\n", user.DisplayName)
			fmt.Printf("User Principal Name: %s\n", user.UserPrincipalName)
			fmt.Printf("Mail Nickname:       %s\n", user.MailNickname)

			return nil
		},
	}

	usersUpdateCmd := &cobra.Command{
		Use:   "update [UPN] [PROPERTY] [VALUE]",
		Short: "Update a user property",
		Long: `Update a user property. Common properties:
  - displayName
  - jobTitle
  - department
  - officeLocation
  - mobilePhone
  - businessPhones (use JSON array format)
  - usageLocation (two-letter country code, e.g., "US")`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			upn := args[0]
			property := args[1]
			value := args[2]

			// Try to parse value as JSON for complex properties
			var parsedValue interface{}
			err := json.Unmarshal([]byte(value), &parsedValue)
			if err != nil {
				// If not valid JSON, treat as string
				parsedValue = value
			}

			properties := map[string]interface{}{
				property: parsedValue,
			}

			err = users.UpdateUser(token, upn, properties)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully updated %s for %s\n", property, upn)
			return nil
		},
	}

	usersDeleteCmd := &cobra.Command{
		Use:   "delete [UPN]",
		Short: "Delete a user",
		Long:  "Delete a user from the tenant. This action can be undone within 30 days by restoring from deleted users.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			upn := args[0]

			// Confirm deletion
			fmt.Printf("⚠ Warning: This will delete user %s\n", upn)
			fmt.Printf("Continue? (yes/no): ")
			var response string
			fmt.Scanln(&response)

			if response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}

			err := users.DeleteUser(token, upn)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully deleted user %s\n", upn)
			return nil
		},
	}

	usersCmd.AddCommand(usersListCmd, usersGetCmd, usersCreateCmd, usersUpdateCmd, usersDeleteCmd)
	rootCmd.AddCommand(usersCmd)
}

// setupLicensesCommands creates the licenses command and its subcommands
func setupLicensesCommands(rootCmd *cobra.Command) {
	licensesCmd := &cobra.Command{
		Use:   "licenses",
		Short: "Manage Microsoft 365 licenses",
	}

	licensesListSkusCmd := &cobra.Command{
		Use:   "list-skus",
		Short: "List all available SKUs in the tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			skus, err := licenses.GetSubscribedSkus(token)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SKU Part Number\tSKU ID\tConsumed Units")
			fmt.Fprintln(w, "---------------\t------\t--------------")
			for _, sku := range skus {
				fmt.Fprintf(w, "%s\t%s\t%d\n", sku.SkuPartNumber, sku.SkuID, sku.ConsumedUnits)
			}
			w.Flush()

			return nil
		},
	}

	licensesGetUserCmd := &cobra.Command{
		Use:   "get [UPN]",
		Short: "Show license details for a specific user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			licenseList, err := licenses.GetUserLicenses(token, args[0])
			if err != nil {
				return err
			}

			if len(licenseList) == 0 {
				fmt.Println("No licenses assigned to this user.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SKU Part Number\tSKU ID")
			fmt.Fprintln(w, "---------------\t------")
			for _, license := range licenseList {
				fmt.Fprintf(w, "%s\t%s\n", license.SkuPartNumber, license.SkuID)
			}
			w.Flush()

			return nil
		},
	}

	licensesAddUserCmd := &cobra.Command{
		Use:   "add-user [UPN] [SKU_ID]",
		Short: "Add a license to a user",
		Long:  "Add one or more licenses to a user by SKU ID. Use 'licenses list-skus' to see available SKU IDs.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			upn := args[0]
			skuIDs := args[1:]

			err := licenses.AssignLicense(token, upn, skuIDs, []string{})
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully added %d license(s) to %s\n", len(skuIDs), upn)
			return nil
		},
	}

	licensesRemoveUserCmd := &cobra.Command{
		Use:   "remove-user [UPN] [SKU_ID]",
		Short: "Remove a license from a user",
		Long:  "Remove one or more licenses from a user by SKU ID. Use 'licenses get [UPN]' to see user's current licenses.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			upn := args[0]
			skuIDs := args[1:]

			err := licenses.AssignLicense(token, upn, []string{}, skuIDs)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully removed %d license(s) from %s\n", len(skuIDs), upn)
			return nil
		},
	}

	licensesGetGroupCmd := &cobra.Command{
		Use:   "get-group [GROUP_ID]",
		Short: "Show license details for a specific group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			licenseList, err := licenses.GetGroupLicenses(token, args[0])
			if err != nil {
				return err
			}

			if len(licenseList) == 0 {
				fmt.Println("No licenses assigned to this group.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "SKU Part Number\tSKU ID")
			fmt.Fprintln(w, "---------------\t------")
			for _, license := range licenseList {
				fmt.Fprintf(w, "%s\t%s\n", license.SkuPartNumber, license.SkuID)
			}
			w.Flush()

			return nil
		},
	}

	licensesAddGroupCmd := &cobra.Command{
		Use:   "add-group [GROUP_ID] [SKU_ID]",
		Short: "Add a license to a group",
		Long:  "Add one or more licenses to a group by SKU ID. Group-based licensing will assign licenses to all members.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := args[0]
			skuIDs := args[1:]

			err := licenses.AssignGroupLicense(token, groupID, skuIDs, []string{})
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully added %d license(s) to group %s\n", len(skuIDs), groupID)
			return nil
		},
	}

	licensesRemoveGroupCmd := &cobra.Command{
		Use:   "remove-group [GROUP_ID] [SKU_ID]",
		Short: "Remove a license from a group",
		Long:  "Remove one or more licenses from a group by SKU ID.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := args[0]
			skuIDs := args[1:]

			err := licenses.AssignGroupLicense(token, groupID, []string{}, skuIDs)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully removed %d license(s) from group %s\n", len(skuIDs), groupID)
			return nil
		},
	}

	licensesCmd.AddCommand(
		licensesListSkusCmd,
		licensesGetUserCmd,
		licensesAddUserCmd,
		licensesRemoveUserCmd,
		licensesGetGroupCmd,
		licensesAddGroupCmd,
		licensesRemoveGroupCmd,
	)

	rootCmd.AddCommand(licensesCmd)
}

// setupGroupsCommands creates the groups command and its subcommands
func setupGroupsCommands(rootCmd *cobra.Command) {
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage Microsoft 365 groups",
	}

	groupsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups in the tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupList, err := groups.ListGroups(token)
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "Display Name\tID\tDescription")
			fmt.Fprintln(w, "------------\t--\t-----------")
			for _, group := range groupList {
				fmt.Fprintf(w, "%s\t%s\t%s\n", group.DisplayName, group.ID, group.Description)
			}
			w.Flush()

			return nil
		},
	}

	groupsGetUserCmd := &cobra.Command{
		Use:   "get [UPN]",
		Short: "Show group memberships for a specific user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupList, err := groups.GetUserGroups(token, args[0])
			if err != nil {
				return err
			}

			if len(groupList) == 0 {
				fmt.Println("User is not a member of any groups.")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "Display Name\tID")
			fmt.Fprintln(w, "------------\t--")
			for _, group := range groupList {
				fmt.Fprintf(w, "%s\t%s\n", group.DisplayName, group.ID)
			}
			w.Flush()

			return nil
		},
	}

	groupsAddUserCmd := &cobra.Command{
		Use:   "add-user [GROUP_ID] [UPN]",
		Short: "Add a user to a group",
		Long:  "Add a user to a group by specifying the group ID and user principal name (UPN).",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := args[0]
			upn := args[1]

			// Get the user ID from UPN
			user, err := users.GetUser(token, upn)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			err = groups.AddMemberToGroup(token, groupID, user.ID)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully added user %s to group %s\n", upn, groupID)
			return nil
		},
	}

	groupsRemoveUserCmd := &cobra.Command{
		Use:   "remove-user [GROUP_ID] [UPN]",
		Short: "Remove a user from a group",
		Long:  "Remove a user from a group by specifying the group ID and user principal name (UPN).",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID := args[0]
			upn := args[1]

			// Get the user ID from UPN
			user, err := users.GetUser(token, upn)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			err = groups.RemoveMemberFromGroup(token, groupID, user.ID)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Successfully removed user %s from group %s\n", upn, groupID)
			return nil
		},
	}

	groupsCmd.AddCommand(groupsListCmd, groupsGetUserCmd, groupsAddUserCmd, groupsRemoveUserCmd)
	rootCmd.AddCommand(groupsCmd)
}
