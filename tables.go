package grist

import (
	"fmt"
	"net/http"
)

type Tables struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	ID      string      `json:"id"`
	Fields  TableFields `json:"fields"`
	Columns []Column    `json:"columns"`
	Records []Record    `json:"records"`
}
type Record struct {
	ID string `json:"id"`
}
type TableFields struct {
	TableRef int  `json:"tableRef"`
	OnDemand bool `json:"onDemand"`
}

type TablePostObj struct {
	Tables []TablePost `json:"tables"`
}

// TablePost represents a table to be created in Grist
// source: https://support.getgrist.com/code/interfaces/DocApiTypes.TablePost/
type TablePost struct {
	ID      string       `json:"id"`
	Records []RecordPost `json:"records,omitempty"`
	Columns []ColumnPost `json:"columns,omitempty"`
}

type RecordPost struct {
	ID     int                   `json:"id,omitempty"`
	Fields map[string]*CellValue `json:"fields"`
}

type ColumnPost struct {
	ID     string               `json:"id,omitempty"`
	Label  string               `json:"label,omitempty"`
	Type   string               `json:"type,omitempty"`
	Fields map[string]CellValue `json:"fields,omitempty"`
}

func pathListTables(docID string) string {
	return pathDescribeDocs(docID) + "/tables"
}

func (d *Doc) ListTables(c *Client) ([]Table, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathListTables(d.ID))
	resp, err := c.GetRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return nil, err
	}

	var t Tables
	if err := handleJSONResponse(resp, &t, http.StatusOK); err != nil {
		return nil, err
	}
	return t.Tables, nil
}

func (d *Doc) CreateTables(c *Client, obj TablePostObj) (*Tables, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathListTables(d.ID))

	jsonBody, err := withJSONBody(obj)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := c.PostRequest(
		endpoint,
		withAuth(c.ApiKey),
		jsonBody,
	)

	if err != nil {
		return nil, err
	}

	var t Tables
	if err := handleJSONResponse(resp, &t, http.StatusOK); err != nil {
		return nil, err
	}

	return &t, nil
}
