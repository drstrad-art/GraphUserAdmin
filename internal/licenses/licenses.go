package licenses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://graph.microsoft.com/v1.0"

// PrepaidUnits represents the prepaid units for a SKU
type PrepaidUnits struct {
	Enabled   int `json:"enabled"`
	Suspended int `json:"suspended"`
	Warning   int `json:"warning"`
}

// SubscribedSku represents a subscribed SKU
type SubscribedSku struct {
	ID            string        `json:"id,omitempty"`
	SkuID         string        `json:"skuId,omitempty"`
	SkuPartNumber string        `json:"skuPartNumber,omitempty"`
	ConsumedUnits int           `json:"consumedUnits,omitempty"`
	PrepaidUnits  PrepaidUnits  `json:"prepaidUnits"`
}

// SkuResponse represents the response when listing SKUs
type SkuResponse struct {
	Value    []SubscribedSku `json:"value"`
	NextLink string          `json:"@odata.nextLink,omitempty"`
}

// LicenseDetail represents a user's license detail
type LicenseDetail struct {
	ID            string `json:"id,omitempty"`
	SkuID         string `json:"skuId,omitempty"`
	SkuPartNumber string `json:"skuPartNumber,omitempty"`
}

// LicenseDetailResponse represents the response when getting user licenses
type LicenseDetailResponse struct {
	Value []LicenseDetail `json:"value"`
}

// AddLicense represents a license to add
type AddLicense struct {
	SkuID string `json:"skuId"`
}

// AssignLicenseRequest represents the request body for assigning licenses
type AssignLicenseRequest struct {
	AddLicenses    []AddLicense `json:"addLicenses"`
	RemoveLicenses []string     `json:"removeLicenses"`
}

// GetSubscribedSkus retrieves all subscribed SKUs for the tenant
func GetSubscribedSkus(accessToken string) ([]SubscribedSku, error) {
	url := fmt.Sprintf("%s/subscribedSkus", baseURL)
	var allSkus []SubscribedSku

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
			return nil, fmt.Errorf("failed to get SKUs: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get SKUs (status %d): %s", resp.StatusCode, string(body))
		}

		var skuResponse SkuResponse
		if err := json.Unmarshal(body, &skuResponse); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allSkus = append(allSkus, skuResponse.Value...)
		url = skuResponse.NextLink
	}

	return allSkus, nil
}

// GetUserLicenses retrieves licenses assigned to a specific user
func GetUserLicenses(accessToken, userPrincipalName string) ([]LicenseDetail, error) {
	url := fmt.Sprintf("%s/users/%s/licenseDetails", baseURL, userPrincipalName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user licenses: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user licenses (status %d): %s", resp.StatusCode, string(body))
	}

	var licenseResponse LicenseDetailResponse
	if err := json.Unmarshal(body, &licenseResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return licenseResponse.Value, nil
}

// AssignLicense adds or removes licenses for a user
func AssignLicense(accessToken, userPrincipalName string, addLicenses []string, removeLicenses []string) error {
	url := fmt.Sprintf("%s/users/%s/assignLicense", baseURL, userPrincipalName)

	// Convert string slices to proper structures
	addLicenseObjs := make([]AddLicense, len(addLicenses))
	for i, skuID := range addLicenses {
		addLicenseObjs[i] = AddLicense{SkuID: skuID}
	}

	assignReq := AssignLicenseRequest{
		AddLicenses:    addLicenseObjs,
		RemoveLicenses: removeLicenses,
	}

	jsonData, err := json.Marshal(assignReq)
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
		return fmt.Errorf("failed to assign license: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Parse error response for better messages
		var errorResponse map[string]interface{}
		if json.Unmarshal(body, &errorResponse) == nil {
			if errorObj, ok := errorResponse["error"].(map[string]interface{}); ok {
				if message, ok := errorObj["message"].(string); ok {
					// Provide helpful context for common errors
					if contains(message, "No available licenses") || contains(message, "license") {
						return fmt.Errorf("license assignment failed: %s\n\nPossible causes:\n  - Not enough available licenses (check with 'gua licenses list-skus')\n  - User doesn't have usageLocation set (use 'gua users update <UPN> usageLocation US')", message)
					}
					if contains(message, "usageLocation") {
						return fmt.Errorf("license assignment failed: %s\n\nUser must have a usageLocation set. Use:\n  gua users update <UPN> usageLocation <country-code>\n  Example: gua users update user@example.com usageLocation US", message)
					}
					return fmt.Errorf("license assignment failed: %s", message)
				}
			}
		}
		return fmt.Errorf("failed to assign license (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if matchAt(s, substr, i) {
			return true
		}
	}
	return false
}

func matchAt(s, substr string, pos int) bool {
	for i := 0; i < len(substr); i++ {
		if toLower(s[pos+i]) != toLower(substr[i]) {
			return false
		}
	}
	return true
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// GetGroupLicenses retrieves licenses assigned to a specific group
func GetGroupLicenses(accessToken, groupID string) ([]LicenseDetail, error) {
	url := fmt.Sprintf("%s/groups/%s/assignedLicenses", baseURL, groupID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get group licenses: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get group licenses (status %d): %s", resp.StatusCode, string(body))
	}

	var licenseResponse LicenseDetailResponse
	if err := json.Unmarshal(body, &licenseResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return licenseResponse.Value, nil
}

// AssignGroupLicense adds or removes licenses for a group
func AssignGroupLicense(accessToken, groupID string, addLicenses []string, removeLicenses []string) error {
	url := fmt.Sprintf("%s/groups/%s/assignLicense", baseURL, groupID)

	// Convert string slices to proper structures
	addLicenseObjs := make([]AddLicense, len(addLicenses))
	for i, skuID := range addLicenses {
		addLicenseObjs[i] = AddLicense{SkuID: skuID}
	}

	assignReq := AssignLicenseRequest{
		AddLicenses:    addLicenseObjs,
		RemoveLicenses: removeLicenses,
	}

	jsonData, err := json.Marshal(assignReq)
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
		return fmt.Errorf("failed to assign group license: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Parse error response for better messages
		var errorResponse map[string]interface{}
		if json.Unmarshal(body, &errorResponse) == nil {
			if errorObj, ok := errorResponse["error"].(map[string]interface{}); ok {
				if message, ok := errorObj["message"].(string); ok {
					// Provide helpful context for common errors
					if contains(message, "No available licenses") || contains(message, "license") {
						return fmt.Errorf("group license assignment failed: %s\n\nPossible causes:\n  - Not enough available licenses (check with 'gua licenses list-skus')\n  - Group members don't have usageLocation set", message)
					}
					return fmt.Errorf("group license assignment failed: %s", message)
				}
			}
		}
		return fmt.Errorf("failed to assign group license (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
