package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/quentinchampenois/go-grist-api"
)

func OrgsExample() {
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

	for _, org := range orgs {
		users, err := org.GetUsersAccess(gc)
		if err != nil {
			fmt.Printf("error describing org: %v: ", err)
			continue
		}

		for _, user := range users {
			fmt.Printf("ID#%d> Email: %s", user.ID, user.Email)
		}
	}
	org := orgs[0]
	// Modify an organization
	org.Name = "Demonstrate"
	err = org.Modify(gc, org.Name)
	if err != nil {
		panic(err)
	}

	// Delete an Organization
	//org.Delete(gc)
	//fmt.Println("Organization deleted")
}
