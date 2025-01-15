package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/openstackclient/identity"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/openstackclient/objectstore"
)

var verbose bool
var containerName string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&containerName, "n", "yesterdaysnews", "container name")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Warn(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	ctx := context.Background()

	identityClient := identity.New(&identity.Config{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		UserID:           os.Getenv("OS_USER_ID"),
		Password:         os.Getenv("OS_PASSWORD"),
	})

	swiftClient := objectstore.New(identityClient, &objectstore.Config{
		Endpoint:   os.Getenv("OS_OBJECTSTORE_URL"),
		ProjectID:  os.Getenv("OS_PROJECT_ID"),
		RegionName: os.Getenv("OS_REGION_NAME"),
	})

	result, err := swiftClient.GetContainer(ctx, containerName)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(data))
}
