package grist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Org struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
	Domain    string     `json:"domain,omitempty"`
	Host      string     `json:"host,omitempty"`
	Access    AccessRole `json:"access"`
}

func ListOrgs(c *Client) ([]Org, error) {
	// Make the GET request
	endpoint := c.ApiEndpoint() + "/orgs"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Add custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Important: Close the response body to avoid resource leaks
	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request %s failed with status: %s", endpoint, resp.Status)
	}
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var orgs []Org
	err = json.Unmarshal(body, &orgs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func DescribeOrg(c *Client, orgId int64) (Org, error) {
	var org Org
	endpoint := c.ApiEndpoint() + "/orgs/" + strconv.FormatInt(orgId, 10)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return org, err
	}
	// Add custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	resp, err := client.Do(req)
	if err != nil {
		return org, err
	}
	defer resp.Body.Close() // Important: Close the response body to avoid resource leaks
	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return org, fmt.Errorf("request %s failed with status: %s", endpoint, resp.Status)
	}
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return org, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &org)
	if err != nil {
		return org, err
	}

	return org, nil
}

func (o *Org) Modify(c *Client, name string) error {
	endpoint := c.ApiEndpoint() + "/orgs/" + strconv.FormatInt(o.ID, 10)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	newName := map[string]string{
		"name": name,
	}

	jsonBody, err := json.MarshalIndent(newName, "", "  ")
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPatch, endpoint, bodyReader)
	if err != nil {
		return err
	}
	// Add custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request %s failed with status: %s", endpoint, resp.Status)
	}

	return nil
}

func (o *Org) Delete(c *Client) error {
	endpoint := c.ApiEndpoint() + "/orgs/" + strconv.FormatInt(o.ID, 10) + "/" + o.Name
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	// Add custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusForbidden:
		return fmt.Errorf("forbidden (status_code: %s)", resp.Status)
	case http.StatusNotFound:
		return fmt.Errorf("not Found (status_code: %s)", resp.Status)
	default:
		fmt.Println("request status:", resp.Status)
	}

	return nil
}

// FIXME: GetUsersAccess depends on deprecated documentation at https://support.getgrist.com/api/#tag/orgs/operation/listOrgAccess
// Open PR to update response documentation
func (o *Org) GetUsersAccess(c *Client) ([]User, error) {
	var (
		accessUsers AccessUser
	)
	endpoint := c.ApiEndpoint() + "/orgs/" + strconv.FormatInt(o.ID, 10) + "/access"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return accessUsers.Users, err
	}
	// Add custom headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	resp, err := client.Do(req)
	if err != nil {
		return accessUsers.Users, err
	}
	defer resp.Body.Close() // Important: Close the response body to avoid resource leaks
	// Check the response status

	fmt.Println(endpoint)
	fmt.Println(o.ID)
	if resp.StatusCode != http.StatusOK {
		return accessUsers.Users, fmt.Errorf("request %s failed with status: %s", endpoint, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return accessUsers.Users, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &accessUsers)
	if err != nil {
		fmt.Println(string(body))
		return accessUsers.Users, err
	}

	return accessUsers.Users, nil
}
