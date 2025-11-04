package grist

import (
	"fmt"
	"net/http"
)

type Doc struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Access    AccessRole `json:"access"`
	IsPinned  bool       `json:"isPinned"`
	UrlID     string     `json:"urlId,omitempty"`
	Workspace Workspace  `json:"workspace,omitempty"`
}

func pathDescribeDocs(docID string) string {
	return "/docs/" + docID
}

// CreateDoc creates a document in the workspace and returns its ID (raw string body).
func (ws *Workspace) CreateDoc(c *Client, name string, isPinned bool) (*string, error) {
	if name == "" {
		return nil, fmt.Errorf("document name cannot be empty")
	}

	endpoint := buildURL(c.ApiEndpoint(), pathWorkspaceDocs(ws.ID))
	bodyOpt, err := withJSONBody(struct {
		Name     string `json:"name"`
		IsPinned bool   `json:"isPinned"`
	}{
		Name:     name,
		IsPinned: isPinned,
	})
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

	b, err := handleRawResponse(resp, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	id := string(b)
	return &id, nil
}

// CreateDoc creates a document in the workspace and returns its ID (raw string body).
// source: https://support.getgrist.com/api/#tag/docs/operation/modifyDoc
// FIXME: Open PR to update response documentation (respond with docID)
func (d *Doc) ModifyMetadataDoc(c *Client, name string, isPinned bool) (*string, error) {
	if name == "" {
		return nil, fmt.Errorf("document name cannot be empty")
	}

	endpoint := buildURL(c.ApiEndpoint(), pathDescribeDocs(d.ID))
	bodyOpt, err := withJSONBody(struct {
		Name     string `json:"name"`
		IsPinned bool   `json:"isPinned"`
	}{
		Name:     name,
		IsPinned: isPinned,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.PatchRequest(
		endpoint,
		withAuth(c.ApiKey),
		bodyOpt,
	)
	if err != nil {
		return nil, err
	}

	b, err := handleRawResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}
	id := string(b)
	return &id, nil
}

// Delete removes a document.
// source: https://support.getgrist.com/api/#tag/docs/operation/deleteDoc
func (d *Doc) Delete(c *Client) error {
	endpoint := buildURL(c.ApiEndpoint(), pathDescribeDocs(d.ID))
	resp, err := c.DeleteRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return err
	}

	// Accept 200 OK
	if err := handleStatus(resp, http.StatusOK); err != nil {
		return err
	}
	return nil
}

// https://support.getgrist.com/api/#tag/docs/operation/importDoc
func (ws *Workspace) ImportExistingDoc(c *Client, url string) (*string, error) {
	return nil, nil
}

// https://support.getgrist.com/api/#tag/docs/operation/describeDoc
func DescribeDoc(c *Client, docID string) (*Doc, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathDescribeDocs(docID))

	resp, err := c.GetRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return nil, err
	}

	var doc Doc
	if err := handleJSONResponse(resp, &doc, http.StatusOK); err != nil {
		return nil, err
	}
	return &doc, nil
}
