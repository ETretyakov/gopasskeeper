package config

import "fmt"

// S3Config is a configuration for S3 storage service.
type S3Config struct {
	Host         string `env:"HOST"              envDefault:"http://127.0.0.1"`
	Port         uint32 `env:"PORT"              envDefault:"9000"`
	AccessKey    string `env:"ACCESS_KEY_ID"     envDefault:"XXXXXXXXXXXXXXXXXXXX"`
	Bucket       string `env:"BUCKET_NAME"       envDefault:"storage-dev"`
	SecretAccess string `env:"SECRET_ACCESS_KEY" envDefault:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`
}

// Endpoint is a method to provide endpoint.
func (s3 *S3Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", s3.Host, s3.Port)
}

// AccessKeyID is a method to provide access key id.
func (s3 *S3Config) AccessKeyID() string {
	return s3.AccessKey
}

// SecretAccessKey is a method to provide secret access key.
func (s3 *S3Config) SecretAccessKey() string {
	return s3.SecretAccess
}

// BucketName is a method to provide bucket name.
func (s3 *S3Config) BucketName() string {
	return s3.Bucket
}
