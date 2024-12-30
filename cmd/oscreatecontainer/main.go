package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/openstackclient/identity"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/openstackclient/objectstore"
)

var verbose bool
var name string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	verbose = true

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	identityClient := identity.New(&identity.Config{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		UserDomainName:   os.Getenv("OS_USER_DOMAIN_NAME"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
	})

	swiftClient := objectstore.New(identityClient, &objectstore.Config{
		Endpoint:   os.Getenv("OS_OBJECTSTORE_URL"),
		ProjectID:  os.Getenv("OS_PROJECT_ID"),
		RegionName: os.Getenv("OS_REGION_NAME"),
	})

	result, err := swiftClient.ListContainers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range result.Containers {
		fmt.Println(container.Name)
	}

	// data, err := swiftClient.Info(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Info(string(data))
}
