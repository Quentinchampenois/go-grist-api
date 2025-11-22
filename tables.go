package grist

import (
	"fmt"
	"net/http"
)

type Tables struct {
	Tables []Table `json:"tables"`
}

// Table represents a table in a Grist document.
// source: https://support.getgrist.com/code/interfaces/DocApiTypes.TablePost/
type Table struct {
	ID      string      `json:"id"`
	Fields  TableFields `json:"fields,omitempty"`
	Columns []Column    `json:"columns,omitempty"`
	Records []Record    `json:"records,omitempty"`
}

// TablesWithColumns represents a table in a Grist document.
// source: https://support.getgrist.com/code/interfaces/DocApiTypes.TablePost/
type TablesWithColumns struct {
	Tables []TableWithColumns `json:"tables"`
}

// TableWithColumns represents a table in a Grist document.
// source: https://support.getgrist.com/code/interfaces/DocApiTypes.TablePost/
type TableWithColumns struct {
	ID      string   `json:"id"`
	Columns []Column `json:"columns"`
}
type Record struct {
	ID     int                   `json:"id,omitempty"`
	Fields map[string]*CellValue `json:"fields"`
}
type TableFields struct {
	TableRef int  `json:"tableRef"`
	OnDemand bool `json:"onDemand"`
}

func pathListTables(docID string) string {
	return pathDescribeDocs(docID) + "/tables"
}

func (d *Doc) ListTables(c *Client) (*Tables, error) {
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
	return &t, nil
}

func (d *Doc) CreateTables(c *Client, obj TablesWithColumns) (*Tables, error) {
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
