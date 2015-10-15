package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BlobstoreClient struct {
	s3Client *s3.S3
}

var errorKeyMissing = errors.New("BlobstoreClient: access_key_id or secret_access_key is missing")
var errorKeysPresent = errors.New("BlobstoreClient: Can't use access_key_id and secret_access_key with env_or_profile credentials_source")

func New(config Config) (*BlobstoreClient, error) {
	transport := *http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: !config.SSLVerifyPeer,
	}

	httpClient := &http.Client{Transport: &transport}

	s3Config := aws.NewConfig().
		WithRegion(config.Region).
		WithS3ForcePathStyle(true).
		// 	// WithEndpoint(config.s3Endpoint()).
		// 	WithLogLevel(aws.LogDebugWithSigning).
		// 	WithDisableSSL(config.UseSSL).
		WithHTTPClient(httpClient)

	switch config.CredentialsSource {
	case "static":
		if config.AccessKeyID == "" || config.SecretAccessKey == "" {
			return nil, errorKeyMissing
		}

		creds := credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, "")
		s3Config = s3Config.WithCredentials(creds)
	case "env_or_profile":
		if config.AccessKeyID != "" || config.SecretAccessKey != "" {
			return nil, errorKeysPresent
		}
	default:
		return nil, fmt.Errorf("Incorrect credentials_source: %s", config.CredentialsSource)
	}

	s3Client := s3.New(s3Config)

	// switch config.Region {
	// case "eu-central-1":
	// 	// use v4 signing
	// case "cn-north-1":
	// 	// use v4 signing
	// default:
	// 	setv2Handlers(s3Client)
	// }

	return &BlobstoreClient{s3Client: s3Client}, nil
}
