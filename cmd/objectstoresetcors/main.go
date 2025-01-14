package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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
		UserDomainName:   os.Getenv("OS_USER_DOMAIN_NAME"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
	})

	swiftClient := objectstore.New(identityClient, &objectstore.Config{
		Endpoint:   os.Getenv("OS_OBJECTSTORE_URL"),
		ProjectID:  os.Getenv("OS_PROJECT_ID"),
		RegionName: os.Getenv("OS_REGION_NAME"),
	})

	headers := map[string]string{
		"X-Container-Meta-Access-Control-Allow-Origin": "*",
	}

	err := swiftClient.SetContainerMetadata(ctx, containerName, headers)
	if err != nil {
		log.Fatal(err)
	}

	corsCheckUrl := fmt.Sprintf("%s/%s", os.Getenv("YN_OBJECTSTORE_URL"), "current.txt")
	err = printCors(ctx, corsCheckUrl)
	if err != nil {
		log.Fatal(err)
	}
}

func printCors(ctx context.Context, url string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodOptions, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("recieved a non-200 status code: %s", resp.Status)
	}

	output := "HEADERS: \n"
	for key, value := range resp.Header {
		output += fmt.Sprintf("%s: %s\n", key, value)
	}
	log.Print(output)

	return nil
}
