package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//
// Config represents the configuration for the s3cli
//
type Config struct {
	AccessKeyID       string `json:"access_key_id"`
	SecretAccessKey   string `json:"secret_access_key"`
	BucketName        string `json:"bucket_name"`
	CredentialsSource string `json:"credentials_source"`
	Host              string `json:"host"`
	Port              int    `json:"port"` // 0 means no custom port
	Region            string `json:"region"`
	SSLVerifyPeer     bool   `json:"ssl_verify_peer"`
	UseSSL            bool   `json:"use_ssl"`
}

func newConfigFromPath(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	config := &Config{
		CredentialsSource: "static",
		Port:              443,
		Region:            "us-east-1",
		SSLVerifyPeer:     true,
		UseSSL:            true,
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	if config.CredentialsSource == "" {
		config.CredentialsSource = "static"
	}

	return config, nil
}
