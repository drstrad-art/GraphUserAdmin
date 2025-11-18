package groups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://graph.microsoft.com/v1.0"

// Group represents a Microsoft 365 group
type Group struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// GroupResponse represents the response when listing groups
type GroupResponse struct {
	Value    []Group `json:"value"`
	NextLink string  `json:"@odata.nextLink,omitempty"`
}

// ListGroups retrieves all groups from Microsoft 365
func ListGroups(accessToken string) ([]Group, error) {
	url := fmt.Sprintf("%s/groups", baseURL)
	var allGroups []Group

	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get groups (status %d): %s", resp.StatusCode, string(body))
		}

		var groupResponse GroupResponse
		if err := json.Unmarshal(body, &groupResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allGroups = append(allGroups, groupResponse.Value...)
		url = groupResponse.NextLink
	}

	return allGroups, nil
}

// GetUserGroups retrieves all groups that a user is a member of
func GetUserGroups(accessToken, userPrincipalName string) ([]Group, error) {
	url := fmt.Sprintf("%s/users/%s/memberOf", baseURL, userPrincipalName)
	var allGroups []Group

	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get user groups: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get user groups (status %d): %s", resp.StatusCode, string(body))
		}

		var groupResponse GroupResponse
		if err := json.Unmarshal(body, &groupResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allGroups = append(allGroups, groupResponse.Value...)
		url = groupResponse.NextLink
	}

	return allGroups, nil
}

// AddMemberToGroup adds a user to a group
func AddMemberToGroup(accessToken, groupID, userID string) error {
	url := fmt.Sprintf("%s/groups/%s/members/$ref", baseURL, groupID)

	// Create request body with the user's directory object ID
	requestBody := map[string]string{
		"@odata.id": fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", userID),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add member to group (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// RemoveMemberFromGroup removes a user from a group
func RemoveMemberFromGroup(accessToken, groupID, userID string) error {
	url := fmt.Sprintf("%s/groups/%s/members/%s/$ref", baseURL, groupID, userID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to remove member from group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove member from group (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
