package grist

import (
	"fmt"
	"net/http"
	"strconv"
)

type Workspace struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Access    AccessRole `json:"access"`
	OrgDomain string     `json:"orgDomain,omitempty"`
	Org       Org        `json:"org,omitempty"`
	Docs      []Doc      `json:"docs,omitempty"`
}

func pathOrgWorkspaces(orgID int64) string {
	return "/orgs/" + strconv.FormatInt(orgID, 10) + "/workspaces"
}

func pathWorkspace(wsID int64) string {
	return "/workspaces/" + strconv.FormatInt(wsID, 10)
}

func pathWorkspaceDocs(wsID int64) string {
	return pathWorkspace(wsID) + "/docs"
}

// ListWorkspaces lists workspaces for an org.
func ListWorkspaces(c *Client, orgId int64) ([]Workspace, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathOrgWorkspaces(orgId))
	fmt.Println(endpoint)
	resp, err := c.GetRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	if err := handleJSONResponse(resp, &workspaces, http.StatusOK); err != nil {
		return nil, err
	}
	return workspaces, nil
}

// CreateWorkspace creates a workspace and returns its ID.
// source: https://support.getgrist.com/api/#tag/workspaces/operation/createWorkspace
func CreateWorkspace(c *Client, orgId int64, name string) (*int64, error) {
	if name == "" {
		return nil, fmt.Errorf("CreateWorkspace: name cannot be empty")
	}

	endpoint := buildURL(c.ApiEndpoint(), pathOrgWorkspaces(orgId))
	bodyOpt, err := withJSONBody(struct {
		Name string `json:"name"`
	}{Name: name})
	if err != nil {
		return nil, err
	}

	resp, err := c.PostRequest(
		endpoint,
		withAuth(c.ApiKey),
		bodyOpt,
	)
	if err != nil {
		return nil, err
	}

	var wsID int64
	err = handleJSONResponse(resp, &wsID, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &wsID, nil
}

// DescribeWorkspace fetches a workspace by ID.
func DescribeWorkspace(c *Client, wsId int64) (*Workspace, error) {
	if wsId <= 0 {
		return nil, fmt.Errorf("invalid wsId: %d", wsId)
	}

	endpoint := buildURL(c.ApiEndpoint(), pathWorkspace(wsId))
	resp, err := c.GetRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return nil, err
	}

	var ws Workspace
	if err := handleJSONResponse(resp, &ws, http.StatusOK); err != nil {
		return nil, err
	}
	return &ws, nil
}

// Modify updates a workspace's name.
func (ws *Workspace) Modify(c *Client, name string) error {
	if ws == nil || ws.ID <= 0 {
		return fmt.Errorf("invalid workspace receiver")
	}
	if name == "" {
		return fmt.Errorf("workspace name cannot be empty")
	}

	endpoint := buildURL(c.ApiEndpoint(), pathWorkspace(ws.ID))
	bodyOpt, err := withJSONBody(struct {
		Name string `json:"name"`
	}{Name: name})
	if err != nil {
		return err
	}

	resp, err := c.PatchRequest(
		endpoint,
		withAuth(c.ApiKey),
		bodyOpt,
	)
	if err != nil {
		return err
	}
	return handleStatus(resp, http.StatusOK)
}

// Delete removes a workspace.
func (ws *Workspace) Delete(c *Client) error {
	if ws == nil || ws.ID <= 0 {
		return fmt.Errorf("invalid workspace receiver")
	}

	endpoint := buildURL(c.ApiEndpoint(), pathWorkspace(ws.ID))
	resp, err := c.DeleteRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return err
	}

	// Accept 200 OK and 204 No Content
	if err := handleStatus(resp, http.StatusOK); err != nil {
		return err
	}
	return nil
}
