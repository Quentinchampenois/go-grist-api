package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/quentinchampenois/go-grist-api"
)

func WorkspacesExample() {
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

	if len(orgs) < 1 {
		panic("No organizations found.")
	}
	org := orgs[0]

	workspaces, err := grist.ListWorkspaces(gc, org.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(workspaces)

	if len(workspaces) > 1 {
		for _, ws := range workspaces {
			err = ws.Delete(gc)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println(workspaces)
	}
	fmt.Println("Creating new workspace...")
	workspace, err := grist.CreateWorkspace(gc, org.ID, "Test Workspace")
	if err != nil {
		panic(err)
	}

	fmt.Println("Workspace created, ID: ", *workspace)

	fmt.Println("Get workspace ID 2...")
	ws, err := grist.DescribeWorkspace(gc, *workspace)
	if err != nil {
		panic(err)
	}

	fmt.Println("Workspace ID 2: ", *ws)

	fmt.Println("Creating new document in workspace...")
	doc, err := ws.CreateDoc(gc, "Test Document", true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document created, ID: ", *doc)

	fmt.Println("Describe document...")
	newDoc, err := grist.DescribeDoc(gc, *doc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Document: ", newDoc)

	fmt.Println("Modifying metadata document...")
	modified, err := newDoc.ModifyMetadataDoc(gc, "Updated name", true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Metadata document modified: ", *modified)

	fmt.Println("Deleting document...")
	err = newDoc.Delete(gc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document deleted.")
	fmt.Println("Deleting workspace...")
	err = ws.Delete(gc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Workspace deleted.")
	fmt.Println("Done.")

}
