package grist

import (
	"fmt"
	"net/http"
)

type Records struct {
	Records []RecordPost `json:"records"`
}

func pathListRecods(docID, tableID string) string {
	return pathDescribeDocs(docID) + "/tables/" + tableID + "/records"
}

// ListRecords lists records for a table.
// https://support.getgrist.com/api/#tag/records/operation/listRecords
func (d *Doc) ListRecords(c *Client, tableID string) (*Records, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathListRecods(d.ID, tableID))
	resp, err := c.GetRequest(
		endpoint,
		withAuth(c.ApiKey),
	)
	if err != nil {
		return nil, err
	}

	var records Records
	err = handleJSONResponse(resp, &records, http.StatusOK)
	if err != nil {
		return nil, err
	}

	return &records, nil
}

// CreateRecords creates new records in a table.
// https://support.getgrist.com/api/#tag/records/operation/addRecords
func (d *Doc) CreateRecords(c *Client, tableID string, obj Records) (*Records, error) {
	endpoint := buildURL(c.ApiEndpoint(), pathListRecods(d.ID, tableID))

	jsonBody, err := withJSONBody(obj)

	fmt.Println("object")
	fmt.Println(obj)
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

	var records Records
	err = handleJSONResponse(resp, &records, http.StatusOK)
	if err != nil {
		return nil, err
	}

	return &records, nil
}
