package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/quentinchampenois/go-grist-api"
)

func DocsExample() {
	fmt.Println("Connecting to Grist...")

	var endpoint string

	if os.Getenv("GRIST_ENV") == "" {
		endpoint = "http://localhost:8484"
	} else {
		endpoint = os.Getenv("GRIST_ENDPOINT")

		if endpoint == "" {
			panic("GRIST_ENDPOINT environment variable not set.")
		}
	}

	apiKey := os.Getenv("GRIST_API_KEY")
	ctx := context.Background()
	gc, err := grist.NewGristClient(ctx, endpoint, apiKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(gc.ApiEndpoint())
	orgs, err := grist.ListOrgs(gc)
	if err != nil {
		panic(err)
	}

	org := orgs[0]
	wsID, err := grist.CreateWorkspace(gc, org.ID, "Doc examples")
	if err != nil {
		fmt.Printf("error creating workspace: %v: ", err)
		return
	}

	ws, err := grist.DescribeWorkspace(gc, *wsID)
	if err != nil {
		fmt.Printf("error describing workspace: %v: ", err)
		return
	}

	docID, err := ws.CreateDoc(gc, "New document", true)
	if err != nil {
		fmt.Printf("error creating new document: %v: ", err)
		return
	}

	doc, err := grist.DescribeDoc(gc, *docID)
	if err != nil {
		fmt.Printf("error describing document: %v: ", err)
	}

	tables, err := doc.ListTables(gc)
	if err != nil {
		fmt.Printf("error listing tables: %v: ", err)
		return
	}

	for _, table := range tables {
		fmt.Println(table.ID)
	}

	gristObj := grist.TablePostObj{
		Tables: []grist.TablePost{
			{
				ID: "Contributors",
				Columns: []grist.ColumnPost{
					{ID: "name", Label: "Name", Type: "Text"},
					{ID: "surname", Label: "Surname", Type: "Text"},
					{ID: "contributions", Label: "Contributions", Type: "Numeric"},
					{ID: "active", Label: "Active", Type: "Boolean"},
				},
			},
		},
	}

	newTables, err := doc.CreateTables(gc, gristObj)
	if err != nil {
		fmt.Printf("error creating tables: %v: ", err)
		return
	}

	if newTables == nil {
		fmt.Println("No tables created.")
		return
	}

	records, err := doc.ListRecords(gc, newTables.Tables[len(newTables.Tables)-1].ID)
	if err != nil {
		fmt.Printf("error listing records: %v: ", err)
		return
	}

	if len(records.Records) < 1 {
		fmt.Println("No records found.")
	}

	name := "Jane Doe"
	surname := "@janedoe"
	contributions := float64(99)
	recordsObj := grist.Records{
		Records: []grist.RecordPost{
			{
				Fields: map[string]*grist.CellValue{
					"name": {
						String: &name,
					},
					"surname": &grist.CellValue{
						String: &surname,
					},
					"contributions": &grist.CellValue{
						Number: &contributions,
					},
				},
			},
		},
	}

	_, err = doc.CreateRecords(gc, newTables.Tables[len(newTables.Tables)-1].ID, recordsObj)
	if err != nil {
		fmt.Printf("error creating records: %v: ", err)
		return
	}

}
