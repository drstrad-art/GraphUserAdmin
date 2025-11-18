package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://graph.microsoft.com/v1.0"

// User represents a Microsoft 365 user
type User struct {
	ID                string `json:"id,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
	Mail              string `json:"mail,omitempty"`
	MailNickname      string `json:"mailNickname,omitempty"`
	AccountEnabled    bool   `json:"accountEnabled,omitempty"`
}

// UserResponse represents the response when listing users
type UserResponse struct {
	Value    []User `json:"value"`
	NextLink string `json:"@odata.nextLink,omitempty"`
}

// PasswordProfile represents password settings for a new user
type PasswordProfile struct {
	ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
	Password                      string `json:"password"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	AccountEnabled    bool            `json:"accountEnabled"`
	DisplayName       string          `json:"displayName"`
	MailNickname      string          `json:"mailNickname"`
	UserPrincipalName string          `json:"userPrincipalName"`
	PasswordProfile   PasswordProfile `json:"passwordProfile"`
}

// ListUsers retrieves all users from Microsoft 365
func ListUsers(accessToken string, filter string) ([]User, error) {
	url := fmt.Sprintf("%s/users", baseURL)
	if filter != "" {
		url += fmt.Sprintf("?$filter=%s", filter)
	}

	var allUsers []User

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
			return nil, fmt.Errorf("failed to get users: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get users (status %d): %s", resp.StatusCode, string(body))
		}

		var userResponse UserResponse
		if err := json.Unmarshal(body, &userResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allUsers = append(allUsers, userResponse.Value...)
		url = userResponse.NextLink
	}

	return allUsers, nil
}

// GetUser retrieves a specific user by UPN
func GetUser(accessToken, userPrincipalName string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, userPrincipalName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user (status %d): %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user in Microsoft 365
func CreateUser(accessToken, displayName, userPrincipalName, mailNickname, password string, forceChange bool) (*User, error) {
	url := fmt.Sprintf("%s/users", baseURL)

	createReq := CreateUserRequest{
		AccountEnabled:    true,
		DisplayName:       displayName,
		MailNickname:      mailNickname,
		UserPrincipalName: userPrincipalName,
		PasswordProfile: PasswordProfile{
			ForceChangePasswordNextSignIn: forceChange,
			Password:                      password,
		},
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create user (status %d): %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// UpdateUser updates properties of an existing user
func UpdateUser(accessToken, userPrincipalName string, properties map[string]interface{}) error {
	url := fmt.Sprintf("%s/users/%s", baseURL, userPrincipalName)

	jsonData, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update user (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteUser deletes a user from Microsoft 365
func DeleteUser(accessToken, userPrincipalName string) error {
	url := fmt.Sprintf("%s/users/%s", baseURL, userPrincipalName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
