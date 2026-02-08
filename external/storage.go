package external

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitStorage(bucketName, accessKey, accessKeySecret, accountId, publicEndpoint string) error {
	// store bucket name for package-level use if needed
	// store base endpoint for constructing public URLs
	baseEndpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, accessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		return err
	}

	// create S3 client with the custom endpoint using service-specific resolver
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Use path-style addressing when necessary. Change to false if using virtual-hosted style.
		o.UsePathStyle = true
		// prefer BaseEndpoint over the deprecated EndpointResolver
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	// assign to package-level variables so other packages can use it (if desired)
	s3Client = client
	bucket = bucketName
	// Normalize pEndpoint by trimming trailing slashes to avoid double slashes in URLs
	pEndpoint = strings.TrimRight(publicEndpoint, "/")

	return nil
}

// exported package-level client and bucket name
var s3Client *s3.Client
var bucket string
var baseEndpoint string
var pEndpoint string

// UploadObject uploads data from r to the configured R2 bucket with the given key and contentType.
// It returns the public URL of the uploaded object.
func UploadObject(ctx context.Context, key string, r io.Reader, contentType string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("s3 client not initialized")
	}

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// Construct public URL
	url := fmt.Sprintf("%s/%s", pEndpoint, key)
	return url, nil
}
