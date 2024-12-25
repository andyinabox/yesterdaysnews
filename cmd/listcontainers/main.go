package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	log.SetLevel(log.DebugLevel)
}

func main() {
	endpoint := os.Getenv("YN_S3_ENDPOINT")
	accessKey := os.Getenv("YN_S3_ACCESS_KEY")
	secretKey := os.Getenv("YN_S3_SECRET_ACCESS_KEY")

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// force path style for openstack compatibility
		o.UsePathStyle = true
		o.BaseEndpoint = &endpoint
	})

	result, err := client.ListBuckets(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range result.Buckets {
		log.Info(*b.Name)
	}

	// Get the first page of results for ListObjectsV2 for a bucket
	// output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
	// 	Bucket: aws.String("my-bucket"),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("first page results:")
	// for _, object := range output.Contents {
	// 	log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	// }

	// endpoint := os.Getenv("YN_S3_ENDPOINT")
	// accessKeyID := os.Getenv("YN_S3_ACCESS_KEY")
	// secretAccessKey := os.Getenv("YN_S3_SECRET_KEY")
	// useSSL := true

	// // Initialize minio client object.
	// log.Debug("initializing client")
	// client, err := minio.New(endpoint, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// 	Region: "us-east-1",
	// })
	// if err != nil {
	// 	log.Fatalf("error intializing client: %s", err)
	// }

	// bucketName := os.Getenv("YN_S3_BUCKET")
	// log.Debugf("checking if bucket %q exists", bucketName)
	// exists, err := client.BucketExists(context.Background(), bucketName)
	// if err != nil {
	// 	log.Fatalf("error checking if bucket exists: %s", err)
	// }
	// log.Infof("bucket exists: %v", exists)

	// log.Debug("listing buckets")
	// buckets, err := client.ListBuckets(context.Background())
	// if err != nil {
	// 	log.Fatalf("error fetching buckets: %s", err)
	// }

	// log.Infof("fetched %d buckets", len(buckets))
	// for _, b := range buckets {
	// 	fmt.Println(b.Name)
	// }
}
