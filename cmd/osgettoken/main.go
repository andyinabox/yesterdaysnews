package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/osidentityclient"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	identityClient := osidentityclient.New(&osidentityclient.Config{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		UserDomainName:   os.Getenv("OS_USER_DOMAIN_NAME"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
	})

	token, err := identityClient.AuthToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}
